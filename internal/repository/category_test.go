package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
)

func TestCategoryRepository_CRUD(t *testing.T) {
	tx := newRepositoryTestTx(t)
	repo := NewCategoryRepository(tx)

	category := &model.Category{
		Name: "test-category-" + uuid.NewString(),
	}

	if err := repo.CreateCategory(category); err != nil {
		t.Fatalf("CreateCategory() error = %v", err)
	}

	gotByID, err := repo.FindCategoryByID(category.ID.String())
	if err != nil {
		t.Fatalf("FindCategoryByID() error = %v", err)
	}
	if gotByID.ID != category.ID {
		t.Fatalf("FindCategoryByID() ID = %s, want %s", gotByID.ID, category.ID)
	}

	all, err := repo.FindAllCategories()
	if err != nil {
		t.Fatalf("FindAllCategories() error = %v", err)
	}
	if len(all) == 0 {
		t.Fatalf("FindAllCategories() returned empty slice")
	}

	category.Name = "updated-" + category.Name
	if err := repo.UpdateCategory(category); err != nil {
		t.Fatalf("UpdateCategory() error = %v", err)
	}

	updated, err := repo.FindCategoryByID(category.ID.String())
	if err != nil {
		t.Fatalf("FindCategoryByID() after update error = %v", err)
	}
	if updated.Name != category.Name {
		t.Fatalf("updated Name = %s, want %s", updated.Name, category.Name)
	}

	if err := repo.DeleteCategory(category.ID.String()); err != nil {
		t.Fatalf("DeleteCategory() error = %v", err)
	}

	if _, err := repo.FindCategoryByID(category.ID.String()); err == nil {
		t.Fatalf("FindCategoryByID() after delete expected error, got nil")
	}
}
