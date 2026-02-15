package model

import "github.com/google/uuid"

type Product struct {
	BaseModel
	Name        string     `gorm:"index;not null" json:"name"`
	SKU         string     `gorm:"uniqueIndex;not null" json:"sku"`
	Description string     `gorm:"type:text" json:"description"`
	Price       float64    `gorm:"not null" json:"price"`
	Stock       int        `gorm:"not null;default:0" json:"stock"`
	CategoryID  *uuid.UUID `gorm:"type:uuid" json:"category_id,omitempty"`
	Category    *Category  `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

type ProductResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
