package repository

import (
	"github.com/lucas-dev/in-memory-db/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)

	FindUserByID(id uint) (*model.User, error)

	UpdateUser(user *model.User) error
	DeleteUser(email string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) RegisterUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUserByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(email string) error {
	return r.db.Where("email = ?", email).Delete(&model.User{}).Error
}
