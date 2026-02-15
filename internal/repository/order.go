package repository

import (
	"github.com/lucas-dev/in-memory-db/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *model.Order) (*model.Order, error)

	GetAllOrders() ([]model.Order, error)
	GetOrderByID(id string) (*model.Order, error)

	UpdateOrder(order *model.Order, id string) (*model.Order, error)
	DeleteOrder(id string) error
}

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
	err := r.db.Preload("User").Preload("Items.Product").First(&result, "id = ?", order.ID).Error
	return &result, err
}

func (r *orderRepository) GetAllOrders() ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("User").Preload("Items.Product").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetOrderByID(id string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("User").Preload("Items.Product").First(&order, "id = ?", id).Error
	return &order, err
}

func (r *orderRepository) UpdateOrder(order *model.Order, id string) (*model.Order, error) {
	return order, r.db.Save(order).
		Where("id = ?", id).
		Error
}

func (r *orderRepository) DeleteOrder(id string) error {
	return r.db.Delete(&model.Order{}, "id = ?", id).Error
}
