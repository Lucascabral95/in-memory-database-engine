package repository

import (
	"github.com/lucas-dev/in-memory-db/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *model.Product) error
	GetAllProducts(limit, offset int) ([]model.Product, int64, error)
	GetProductByID(id string) (*model.Product, error)

	UpdateProduct(product *model.Product, id string) error
	DeleteProduct(id string) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetAllProducts(limit, offset int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	if err := r.db.Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) GetProductByID(id string) (*model.Product, error) {
	var product model.Product
	return &product, r.db.First(&product, "id = ?", id).Error
}

func (r *productRepository) DeleteProduct(id string) error {
	return r.db.Delete(&model.Product{}, "id = ?", id).Error
}

func (r *productRepository) UpdateProduct(product *model.Product, id string) error {
	return r.db.Model(&model.Product{}).
		Where("id = ?", id).
		Select("Name", "Description", "Price", "Stock", "CategoryID").
		Updates(product).Error
}
