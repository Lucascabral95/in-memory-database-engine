package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StockMovementReason representa el motivo de un movimiento de stock
type StockMovementReason string

const (
	StockMovementReasonSale       StockMovementReason = "SALE"
	StockMovementReasonRestock    StockMovementReason = "RESTOCK"
	StockMovementReasonAdjustment StockMovementReason = "ADJUSTMENT"
)

// StockMovement representa un movimiento de stock
type StockMovement struct {
	ID        uuid.UUID           `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID uuid.UUID           `gorm:"type:uuid;not null;index" json:"product_id"`
	Product   Product             `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	Quantity  int                 `gorm:"not null" json:"quantity"`
	Reason    StockMovementReason `gorm:"size:255" json:"reason"`
	CreatedAt time.Time           `json:"created_at"`
}

func (sm *StockMovement) BeforeCreate(tx *gorm.DB) error {
	if sm.ID == uuid.Nil {
		sm.ID = uuid.New()
	}
	return nil
}
