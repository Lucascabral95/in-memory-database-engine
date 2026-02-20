package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
)

func TestProductRepository_CRUDAndPagination(t *testing.T) {
	tx := newRepositoryTestTx(t)
	repo := NewProductRepository(tx)

	product := &model.Product{
		Name:        "test-product-" + uuid.NewString(),
		SKU:         "SKU-" + uuid.NewString(),
		Description: "product for repository test",
		Price:       99.90,
		Stock:       10,
	}

	if err := repo.CreateProduct(product); err != nil {
		t.Fatalf("CreateProduct() error = %v", err)
	}

	gotByID, err := repo.GetProductByID(product.ID.String())
	if err != nil {
		t.Fatalf("GetProductByID() error = %v", err)
	}
	if gotByID.ID != product.ID {
		t.Fatalf("GetProductByID() ID = %s, want %s", gotByID.ID, product.ID)
	}

	list, total, err := repo.GetAllProducts(10, 0)
	if err != nil {
		t.Fatalf("GetAllProducts() error = %v", err)
	}
	if total == 0 || len(list) == 0 {
		t.Fatalf("GetAllProducts() expected data, got total=%d len=%d", total, len(list))
	}

	product.Name = "updated-" + product.Name
	product.Stock = 7
	if err := repo.UpdateProduct(product, product.ID.String()); err != nil {
		t.Fatalf("UpdateProduct() error = %v", err)
	}

	updated, err := repo.GetProductByID(product.ID.String())
	if err != nil {
		t.Fatalf("GetProductByID() after update error = %v", err)
	}
	if updated.Name != product.Name || updated.Stock != product.Stock {
		t.Fatalf("updated product mismatch: got name=%s stock=%d", updated.Name, updated.Stock)
	}

	if err := repo.DeleteProduct(product.ID.String()); err != nil {
		t.Fatalf("DeleteProduct() error = %v", err)
	}
	if _, err := repo.GetProductByID(product.ID.String()); err == nil {
		t.Fatalf("GetProductByID() after delete expected error, got nil")
	}
}

func TestProductRepository_GetAllProducts_RespectsLimitAndOffset(t *testing.T) {
	tx := newRepositoryTestTx(t)
	repo := NewProductRepository(tx)

	for i := 0; i < 3; i++ {
		product := &model.Product{
			Name:        "paged-product-" + uuid.NewString(),
			SKU:         "SKU-" + uuid.NewString(),
			Description: "pagination test product",
			Price:       50 + float64(i),
			Stock:       5 + i,
		}
		if err := repo.CreateProduct(product); err != nil {
			t.Fatalf("CreateProduct() error = %v", err)
		}
	}

	page, total, err := repo.GetAllProducts(1, 1)
	if err != nil {
		t.Fatalf("GetAllProducts() error = %v", err)
	}
	if total < 3 {
		t.Fatalf("total = %d, want >= 3", total)
	}
	if len(page) != 1 {
		t.Fatalf("len(page) = %d, want %d", len(page), 1)
	}
}

func TestProductRepository_GetProductByID_NotFound(t *testing.T) {
	tx := newRepositoryTestTx(t)
	repo := NewProductRepository(tx)

	_, err := repo.GetProductByID(uuid.NewString())
	if err == nil {
		t.Fatalf("GetProductByID() expected error for non-existing product")
	}
}
