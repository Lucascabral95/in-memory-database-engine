package model

type Category struct {
	BaseModel
	Name     string    `gorm:"uniqueIndex;not null" json:"name"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
