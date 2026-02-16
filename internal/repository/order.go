package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) (*model.Order, error)

	GetAllOrders() ([]model.Order, error)
	GetOrdersByUserID(userID string) ([]model.Order, error)
	GetOrderByID(id string) (*model.Order, error)

	UpdateStatusOrder(order *model.UpdateOrderStatusRequest, id string) error
	// DeleteOrder(id string) error
	OrderUpdatePay(id string) error
}

var (
	ErrOrderNotPending   = errors.New("la orden no está en estado PENDING")
	ErrInsufficientStock = errors.New("stock insuficiente para completar el pago")
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *model.Order) (*model.Order, error) {
	if err := r.db.Create(order).Error; err != nil {
		return nil, err
	}

	var result model.Order
	err := r.db.
		Preload("User").
		Preload("Items.Product").
		First(&result, "id = ?", order.ID).
		Error
	return &result, err
}

func (r *orderRepository) GetAllOrders() ([]model.Order, error) {
	var orders []model.Order
	err := r.db.
		Preload("User").
		Preload("Items.Product").
		Find(&orders).
		Error
	return orders, err
}

func (r *orderRepository) GetOrdersByUserID(userID string) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.
		Preload("User").
		Preload("Items.Product").
		Where("user_id = ?", userID).
		Find(&orders).
		Error
	return orders, err
}

func (r *orderRepository) GetOrderByID(id string) (*model.Order, error) {
	var order model.Order
	err := r.db.
		Preload("User").
		Preload("Items.Product").
		First(&order, "id = ?", id).
		Error
	return &order, err
}

// func (r *orderRepository) DeleteOrder(id string) error {
// 	return r.db.Delete(&model.Order{}, "id = ?", id).Error
// }

func (r *orderRepository) UpdateStatusOrder(order *model.UpdateOrderStatusRequest, id string) error {
	return r.db.Model(&model.Order{}).
		Where("id = ?", id).
		Update("status", order.Status).
		Error
}

func (r *orderRepository) OrderUpdatePay(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.Preload("Items").First(&order, "id = ?", id).Error; err != nil {
			return err
		}

		if order.Status != model.OrderStatusPending {
			return ErrOrderNotPending
		}

		requiredByProduct := make(map[uuid.UUID]int)
		for _, item := range order.Items {
			requiredByProduct[item.ProductID] += item.Quantity
		}

		for productID, requiredQty := range requiredByProduct {
			var product model.Product
			if err := tx.
				Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&product, "id = ?", productID).Error; err != nil {
				return err
			}

			if product.Stock < requiredQty {
				return ErrInsufficientStock
			}

			if err := tx.Model(&model.Product{}).
				Where("id = ?", productID).
				Update("stock", gorm.Expr("stock - ?", requiredQty)).
				Error; err != nil {
				return err
			}

			movement := model.StockMovement{
				ProductID: productID,
				Quantity:  -requiredQty,
				Reason:    model.StockMovementReasonSale,
			}

			if err := tx.Create(&movement).Error; err != nil {
				return err
			}
		}

		return tx.Model(&model.Order{}).
			Where("id = ?", id).
			Update("status", model.OrderStatusPaid).
			Error
	})
}
