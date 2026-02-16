package repository

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
)

func TestOrderRepository_CreateAndGetByID(t *testing.T) {
	tx := newRepositoryTestTx(t)
	repo := NewOrderRepository(tx)

	user := seedUser(t, tx)
	product := seedProduct(t, tx, 10)

	order := &model.Order{
		UserID:      user.ID,
		Status:      model.OrderStatusPending,
		TotalAmount: 120,
		Items: []model.OrderItem{
			{
				ProductID:     product.ID,
				Quantity:      2,
				PriceAtMoment: 60,
			},
		},
	}

	created, err := repo.CreateOrder(order)
	if err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}
	if created.ID == uuid.Nil {
		t.Fatalf("CreateOrder() returned nil UUID")
	}

	got, err := repo.GetOrderByID(created.ID.String())
	if err != nil {
		t.Fatalf("GetOrderByID() error = %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("GetOrderByID() ID = %s, want %s", got.ID, created.ID)
	}
}

func TestOrderRepository_OrderUpdatePay_Success(t *testing.T) {
	tx := newRepositoryTestTx(t)
	repo := NewOrderRepository(tx)

	user := seedUser(t, tx)
	product := seedProduct(t, tx, 5)

	order := seedOrder(t, tx, user.ID, product.ID, 3, 100)

	if err := repo.OrderUpdatePay(order.ID.String()); err != nil {
		t.Fatalf("OrderUpdatePay() error = %v", err)
	}

	updatedOrder, err := repo.GetOrderByID(order.ID.String())
	if err != nil {
		t.Fatalf("GetOrderByID() after pay error = %v", err)
	}
	if updatedOrder.Status != model.OrderStatusPaid {
		t.Fatalf("order status = %s, want %s", updatedOrder.Status, model.OrderStatusPaid)
	}

	var updatedProduct model.Product
	if err := tx.First(&updatedProduct, "id = ?", product.ID).Error; err != nil {
		t.Fatalf("loading product after pay error = %v", err)
	}
	if updatedProduct.Stock != 2 {
		t.Fatalf("product stock = %d, want %d", updatedProduct.Stock, 2)
	}

	var movements []model.StockMovement
	if err := tx.Where("product_id = ?", product.ID).Find(&movements).Error; err != nil {
		t.Fatalf("loading stock movements error = %v", err)
	}
	if len(movements) != 1 {
		t.Fatalf("stock movements = %d, want %d", len(movements), 1)
	}
	if movements[0].Reason != model.StockMovementReasonSale {
		t.Fatalf("stock movement reason = %s, want %s", movements[0].Reason, model.StockMovementReasonSale)
	}
}

func TestOrderRepository_OrderUpdatePay_InsufficientStock(t *testing.T) {
	tx := newRepositoryTestTx(t)
	repo := NewOrderRepository(tx)

	user := seedUser(t, tx)
	product := seedProduct(t, tx, 1)
	order := seedOrder(t, tx, user.ID, product.ID, 2, 50)

	err := repo.OrderUpdatePay(order.ID.String())
	if !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("OrderUpdatePay() error = %v, want %v", err, ErrInsufficientStock)
	}

	stillPending, err := repo.GetOrderByID(order.ID.String())
	if err != nil {
		t.Fatalf("GetOrderByID() after failed pay error = %v", err)
	}
	if stillPending.Status != model.OrderStatusPending {
		t.Fatalf("order status = %s, want %s", stillPending.Status, model.OrderStatusPending)
	}
}

func TestOrderRepository_OrderUpdatePay_NotPending(t *testing.T) {
	tx := newRepositoryTestTx(t)
	repo := NewOrderRepository(tx)

	user := seedUser(t, tx)
	product := seedProduct(t, tx, 5)
	order := seedOrder(t, tx, user.ID, product.ID, 1, 100)

	if err := tx.Model(&model.Order{}).
		Where("id = ?", order.ID).
		Update("status", model.OrderStatusPaid).
		Error; err != nil {
		t.Fatalf("preparing paid status error = %v", err)
	}

	err := repo.OrderUpdatePay(order.ID.String())
	if !errors.Is(err, ErrOrderNotPending) {
		t.Fatalf("OrderUpdatePay() error = %v, want %v", err, ErrOrderNotPending)
	}
}

func seedUser(t *testing.T, tx repositoryTestDB) model.User {
	t.Helper()

	user := model.User{
		Email:     "order-user-" + uuid.NewString() + "@example.com",
		Password:  "hashed-password",
		FirstName: "Order",
		LastName:  "Tester",
	}
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("seed user error = %v", err)
	}
	return user
}

func seedProduct(t *testing.T, tx repositoryTestDB, stock int) model.Product {
	t.Helper()

	product := model.Product{
		Name:        "order-product-" + uuid.NewString(),
		SKU:         "SKU-" + uuid.NewString(),
		Description: "product for order repo tests",
		Price:       100,
		Stock:       stock,
	}
	if err := tx.Create(&product).Error; err != nil {
		t.Fatalf("seed product error = %v", err)
	}
	return product
}

func seedOrder(
	t *testing.T,
	tx repositoryTestDB,
	userID uuid.UUID,
	productID uuid.UUID,
	quantity int,
	price float64,
) model.Order {
	t.Helper()

	order := model.Order{
		UserID:      userID,
		Status:      model.OrderStatusPending,
		TotalAmount: float64(quantity) * price,
		Items: []model.OrderItem{
			{
				ProductID:     productID,
				Quantity:      quantity,
				PriceAtMoment: price,
			},
		},
	}

	if err := tx.Create(&order).Error; err != nil {
		t.Fatalf("seed order error = %v", err)
	}

	return order
}
