package model

// Category representa una categoría de productos
type Category struct {
	BaseModel
	Name     string    `gorm:"uniqueIndex;not null" json:"name"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
