package repository

import (
	"github.com/lucas-dev/in-memory-db/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *model.Category) error

	FindAllCategories() ([]model.Category, error)
	FindCategoryByID(id string) (*model.Category, error)

	UpdateCategory(category *model.Category) error
	DeleteCategory(id string) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindAllCategories() ([]model.Category, error) {
	var categories []model.Category
	return categories, r.db.Find(&categories).Error
}

func (r *categoryRepository) FindCategoryByID(id string) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) DeleteCategory(id string) error {
	err := r.db.Delete(&model.Category{}, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) UpdateCategory(category *model.Category) error {
	return r.db.Save(category).Error
}
