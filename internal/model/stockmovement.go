package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockMovement struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Reason    string    `gorm:"size:255" json:"reason"` // "SALE", "RESTOCK", "ADJUSTMENT"
	CreatedAt time.Time `json:"created_at"`
}

func (sm *StockMovement) BeforeCreate(tx *gorm.DB) error {
	if sm.ID == uuid.Nil {
		sm.ID = uuid.New()
	}
	return nil
}
