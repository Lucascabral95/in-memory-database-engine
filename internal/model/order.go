package model

import "github.com/google/uuid"

// Order representa una orden de compra
type Order struct {
	BaseModel
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`

	TotalAmount float64     `gorm:"not null" json:"total_amount"`
	Status      OrderStatus `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	Items       []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

// OrderStatus representa el estado de una orden
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// UpdateOrderStatusRequest representa la solicitud para actualizar el estado de una orden
type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=PENDING PAID SHIPPED CANCELLED"`
}
