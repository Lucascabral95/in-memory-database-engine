package model

import "github.com/google/uuid"

// OrderItem representa un producto en una orden
type OrderItem struct {
	BaseModel
	OrderID       uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Product       Product   `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	PriceAtMoment float64   `gorm:"not null" json:"price_at_moment"`
}
