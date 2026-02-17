package database

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/pkg/utils"
	"gorm.io/gorm"
)

const (
	defaultSeedUsers      = 24
	defaultSeedCategories = 8
	defaultSeedProducts   = 64
	defaultSeedOrders     = 50
)

type SeedOptions struct {
	Users      int
	Categories int
	Products   int
	Orders     int
	Force      bool
}

func (o SeedOptions) normalize() SeedOptions {
	if o.Users <= 0 {
		o.Users = defaultSeedUsers
	}
	if o.Categories <= 0 {
		o.Categories = defaultSeedCategories
	}
	if o.Products <= 0 {
		o.Products = defaultSeedProducts
	}
	if o.Orders <= 0 {
		o.Orders = defaultSeedOrders
	}
	return o
}

func SeedDatabase(db *gorm.DB, opts SeedOptions) error {
	opts = opts.normalize()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	return db.Transaction(func(tx *gorm.DB) error {
		if opts.Force {
			if err := purgeSeedTables(tx); err != nil {
				return err
			}
		} else {
			if err := ensureEmptyDB(tx); err != nil {
				return err
			}
		}

		categories, err := seedCategories(tx, opts.Categories)
		if err != nil {
			return err
		}

		users, err := seedUsers(tx, opts.Users)
		if err != nil {
			return err
		}

		products, stockByProduct, err := seedProducts(tx, opts.Products, categories, rng)
		if err != nil {
			return err
		}

		if err := seedOrders(tx, opts.Orders, users, products, stockByProduct, rng); err != nil {
			return err
		}

		return nil
	})
}

func purgeSeedTables(tx *gorm.DB) error {
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.StockMovement{}).Error; err != nil {
		return fmt.Errorf("purge stock movements: %w", err)
	}
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.OrderItem{}).Error; err != nil {
		return fmt.Errorf("purge order items: %w", err)
	}
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.Order{}).Error; err != nil {
		return fmt.Errorf("purge orders: %w", err)
	}
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.Product{}).Error; err != nil {
		return fmt.Errorf("purge products: %w", err)
	}
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.Category{}).Error; err != nil {
		return fmt.Errorf("purge categories: %w", err)
	}
	if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&model.User{}).Error; err != nil {
		return fmt.Errorf("purge users: %w", err)
	}
	return nil
}

func ensureEmptyDB(db *gorm.DB) error {
	type counter struct {
		model interface{}
		name  string
	}

	counters := []counter{
		{model: &model.User{}, name: "users"},
		{model: &model.Category{}, name: "categories"},
		{model: &model.Product{}, name: "products"},
		{model: &model.Order{}, name: "orders"},
	}

	for _, c := range counters {
		var count int64
		if err := db.Model(c.model).Count(&count).Error; err != nil {
			return fmt.Errorf("count %s: %w", c.name, err)
		}
		if count > 0 {
			return fmt.Errorf("seed cancelled: table %s already has data (%d rows)", c.name, count)
		}
	}

	return nil
}

func seedCategories(tx *gorm.DB, total int) ([]model.Category, error) {
	baseNames := []string{
		"Electronics",
		"Gaming",
		"Home",
		"Fashion",
		"Sports",
		"Books",
		"Beauty",
		"Toys",
	}

	categories := make([]model.Category, 0, total)
	for i := 0; i < total; i++ {
		name := baseNames[i%len(baseNames)]
		if i >= len(baseNames) {
			name = fmt.Sprintf("%s %d", name, i+1)
		}
		categories = append(categories, model.Category{Name: name})
	}

	if err := tx.Create(&categories).Error; err != nil {
		return nil, fmt.Errorf("create categories: %w", err)
	}

	return categories, nil
}

func seedUsers(tx *gorm.DB, total int) ([]model.User, error) {
	password, err := utils.HashPassword("Password123!")
	if err != nil {
		return nil, fmt.Errorf("hash seed password: %w", err)
	}

	users := make([]model.User, 0, total)
	for i := 1; i <= total; i++ {
		users = append(users, model.User{
			Email:     fmt.Sprintf("seed.user%02d@example.com", i),
			Password:  password,
			FirstName: fmt.Sprintf("User%02d", i),
			LastName:  "Seed",
		})
	}

	if err := tx.Create(&users).Error; err != nil {
		return nil, fmt.Errorf("create users: %w", err)
	}

	return users, nil
}

func seedProducts(tx *gorm.DB, total int, categories []model.Category, rng *rand.Rand) ([]model.Product, map[uuid.UUID]int, error) {
	adjectives := []string{"Smart", "Pro", "Ultra", "Classic", "Prime", "Compact", "Max", "Air", "Eco", "Elite"}
	objects := []string{"Mouse", "Keyboard", "Headset", "Monitor", "Chair", "Camera", "Speaker", "Watch", "Backpack", "Lamp"}

	products := make([]model.Product, 0, total)
	stockByProduct := make(map[uuid.UUID]int, total)

	for i := 1; i <= total; i++ {
		category := categories[i%len(categories)]
		name := fmt.Sprintf("%s %s %d", adjectives[i%len(adjectives)], objects[i%len(objects)], i)
		price := 10 + rng.Float64()*490
		stock := rng.Intn(121) + 80
		sku := fmt.Sprintf("SKU-%03d-%s", i, uuid.NewString()[:8])

		product := model.Product{
			Name:        name,
			SKU:         sku,
			Description: fmt.Sprintf("%s for category %s", name, category.Name),
			Price:       round2(price),
			Stock:       stock,
			CategoryID:  &category.ID,
		}

		products = append(products, product)
	}

	if err := tx.Create(&products).Error; err != nil {
		return nil, nil, fmt.Errorf("create products: %w", err)
	}

	for _, p := range products {
		stockByProduct[p.ID] = p.Stock
	}

	return products, stockByProduct, nil
}

func seedOrders(tx *gorm.DB, total int, users []model.User, products []model.Product, stockByProduct map[uuid.UUID]int, rng *rand.Rand) error {
	for i := 0; i < total; i++ {
		user := users[rng.Intn(len(users))]
		itemsCount := rng.Intn(4) + 1 // 1..4
		status := randomOrderStatus(rng)

		selectedProducts := pickUniqueProducts(products, itemsCount, rng)
		orderItems := make([]model.OrderItem, 0, len(selectedProducts))
		requiredQtyByProduct := make(map[uuid.UUID]int, len(selectedProducts))
		totalAmount := 0.0

		for _, p := range selectedProducts {
			qty := rng.Intn(3) + 1
			requiredQtyByProduct[p.ID] += qty
			totalAmount += p.Price * float64(qty)
			orderItems = append(orderItems, model.OrderItem{
				ProductID:     p.ID,
				Quantity:      qty,
				PriceAtMoment: p.Price,
			})
		}

		if (status == model.OrderStatusPaid || status == model.OrderStatusShipped) && !canFulfill(requiredQtyByProduct, stockByProduct) {
			status = model.OrderStatusCancelled
		}

		order := model.Order{
			UserID:      user.ID,
			TotalAmount: round2(totalAmount),
			Status:      status,
		}

		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("create order: %w", err)
		}

		for idx := range orderItems {
			orderItems[idx].OrderID = order.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return fmt.Errorf("create order items: %w", err)
		}

		if status == model.OrderStatusPaid || status == model.OrderStatusShipped {
			for productID, qty := range requiredQtyByProduct {
				stockByProduct[productID] -= qty

				if err := tx.Model(&model.Product{}).
					Where("id = ?", productID).
					Update("stock", gorm.Expr("stock - ?", qty)).
					Error; err != nil {
					return fmt.Errorf("update product stock: %w", err)
				}

				movement := model.StockMovement{
					ProductID: productID,
					Quantity:  -qty,
					Reason:    model.StockMovementReasonSale,
				}
				if err := tx.Create(&movement).Error; err != nil {
					return fmt.Errorf("create stock movement: %w", err)
				}
			}
		}
	}

	return nil
}

func randomOrderStatus(rng *rand.Rand) model.OrderStatus {
	roll := rng.Intn(100)
	switch {
	case roll < 50:
		return model.OrderStatusPaid
	case roll < 70:
		return model.OrderStatusPending
	case roll < 90:
		return model.OrderStatusShipped
	default:
		return model.OrderStatusCancelled
	}
}

func pickUniqueProducts(products []model.Product, total int, rng *rand.Rand) []model.Product {
	if total >= len(products) {
		return products
	}

	used := make(map[uuid.UUID]struct{}, total)
	selected := make([]model.Product, 0, total)
	for len(selected) < total {
		p := products[rng.Intn(len(products))]
		if _, ok := used[p.ID]; ok {
			continue
		}
		used[p.ID] = struct{}{}
		selected = append(selected, p)
	}
	return selected
}

func canFulfill(required map[uuid.UUID]int, stockByProduct map[uuid.UUID]int) bool {
	for productID, qty := range required {
		if stockByProduct[productID] < qty {
			return false
		}
	}
	return true
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}
