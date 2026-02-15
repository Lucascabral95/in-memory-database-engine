package model

type User struct {
	BaseModel
	Email     string  `gorm:"uniqueIndex;not null" json:"email"`
	Password  string  `gorm:"size:60;not null" json:"-"`
	FirstName string  `gorm:"size:100" json:"first_name"`
	LastName  string  `gorm:"size:100" json:"last_name"`
	Orders    []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

// DTOs de users
type UserCreateRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
