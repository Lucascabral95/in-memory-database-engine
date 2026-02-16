package model

import "github.com/google/uuid"

type CartItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
}

type Cart struct {
	UserID uuid.UUID  `json:"user_id"`
	Items  []CartItem `json:"items"`
}
