package model

import "github.com/google/uuid"

// CartItem representa un item en el carrito de compras
type CartItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
}

// Cart representa el carrito de compras de un usuario
type Cart struct {
	UserID uuid.UUID  `json:"user_id"`
	Items  []CartItem `json:"items"`
}
