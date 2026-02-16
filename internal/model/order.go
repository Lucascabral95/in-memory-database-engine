package model

import "github.com/google/uuid"

type Order struct {
	BaseModel
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`

	TotalAmount float64     `gorm:"not null" json:"total_amount"`
	Status      OrderStatus `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	Items       []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=PENDING PAID SHIPPED CANCELLED"`
}
