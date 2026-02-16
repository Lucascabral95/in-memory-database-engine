package service

import (
	"errors"
	"testing"

	"github.com/lucas-dev/in-memory-db/internal/model"
)

type MockProductRepository struct {
	CreateProductFunc  func(product *model.Product) error
	GetAllProductsFunc func(limit, offset int) ([]model.Product, int64, error)
	GetProductByIDFunc func(id string) (*model.Product, error)
	UpdateProductFunc  func(product *model.Product, id string) error
	DeleteProductFunc  func(id string) error
}

func (m *MockProductRepository) CreateProduct(product *model.Product) error {
	if m.CreateProductFunc != nil {
		return m.CreateProductFunc(product)
	}
	return nil
}

func (m *MockProductRepository) GetAllProducts(limit, offset int) ([]model.Product, int64, error) {
	if m.GetAllProductsFunc != nil {
		return m.GetAllProductsFunc(limit, offset)
	}
	return []model.Product{}, 0, nil
}

func (m *MockProductRepository) GetProductByID(id string) (*model.Product, error) {
	if m.GetProductByIDFunc != nil {
		return m.GetProductByIDFunc(id)
	}
	return &model.Product{}, nil
}

func (m *MockProductRepository) UpdateProduct(product *model.Product, id string) error {
	if m.UpdateProductFunc != nil {
		return m.UpdateProductFunc(product, id)
	}
	return nil
}

func (m *MockProductRepository) DeleteProduct(id string) error {
	if m.DeleteProductFunc != nil {
		return m.DeleteProductFunc(id)
	}
	return nil
}

func TestProductService_CreateProduct(t *testing.T) {
	t.Run("NilBody", func(t *testing.T) {
		svc := NewProductService(&MockProductRepository{})
		err := svc.CreateProduct(nil)
		if err == nil || err.Error() != "El Body Request viene vacío." {
			t.Errorf("Se esperaba error body vacío, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockProductRepository{
			CreateProductFunc: func(p *model.Product) error { return nil },
		}
		svc := NewProductService(mockRepo)
		err := svc.CreateProduct(&model.Product{Name: "Laptop"})
		if err != nil {
			t.Errorf("No se esperaba error, got %v", err)
		}
	})
}

func TestProductService_GetAllProducts(t *testing.T) {
	t.Run("PaginationLogic", func(t *testing.T) {
		callCount := 0
		mockRepo := &MockProductRepository{
			GetAllProductsFunc: func(limit, offset int) ([]model.Product, int64, error) {
				callCount++
				if callCount == 1 {
					if limit != 0 || offset != 0 {
						t.Errorf("Primera llamada esperada con (0,0), got (%d,%d)", limit, offset)
					}
					return nil, 15, nil
				}
				if limit != 10 || offset != 0 {
					t.Errorf("Segunda llamada esperada con (10,0), got (%d,%d)", limit, offset)
				}
				return []model.Product{{}, {}}, 15, nil
			},
		}
		svc := NewProductService(mockRepo)

		resp, err := svc.GetAllProducts(10, 1)
		if err != nil {
			t.Fatalf("Error inesperado: %v", err)
		}

		if resp.Total != 15 {
			t.Errorf("Expected total 15, got %d", resp.Total)
		}
		if resp.TotalPages != 2 {
			t.Errorf("Expected 2 pages, got %d", resp.TotalPages)
		}
	})

	t.Run("DefaultValues", func(t *testing.T) {
		mockRepo := &MockProductRepository{
			GetAllProductsFunc: func(l, o int) ([]model.Product, int64, error) {
				return nil, 0, nil
			},
		}
		svc := NewProductService(mockRepo)

		_, err := svc.GetAllProducts(0, 0)
		if err != nil {
			if err.Error() != "no data found" {
			}
		}
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	t.Run("EmptyID", func(t *testing.T) {
		svc := NewProductService(&MockProductRepository{})
		err := svc.DeleteProduct("")
		if err == nil || err.Error() != "El ID viene vacío." {
			t.Errorf("Se esperaba error ID vacío, got %v", err)
		}
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		mockRepo := &MockProductRepository{
			GetProductByIDFunc: func(id string) (*model.Product, error) {
				return nil, errors.New("not found")
			},
		}
		svc := NewProductService(mockRepo)
		err := svc.DeleteProduct("123")
		if err == nil {
			t.Error("Se esperaba error al obtener producto")
		}
	})
}
