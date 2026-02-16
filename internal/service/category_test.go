package service

import (
	"errors"
	"testing"

	"github.com/lucas-dev/in-memory-db/internal/model"
)

type MockCategoryRepository struct {
	CreateCategoryFunc    func(category *model.Category) error
	FindAllCategoriesFunc func() ([]model.Category, error)
	FindCategoryByIDFunc  func(id string) (*model.Category, error)
	UpdateCategoryFunc    func(category *model.Category) error
	DeleteCategoryFunc    func(id string) error
}

func (m *MockCategoryRepository) CreateCategory(category *model.Category) error {
	if m.CreateCategoryFunc != nil {
		return m.CreateCategoryFunc(category)
	}
	return nil
}

func (m *MockCategoryRepository) FindAllCategories() ([]model.Category, error) {
	if m.FindAllCategoriesFunc != nil {
		return m.FindAllCategoriesFunc()
	}
	return []model.Category{}, nil
}

func (m *MockCategoryRepository) FindCategoryByID(id string) (*model.Category, error) {
	if m.FindCategoryByIDFunc != nil {
		return m.FindCategoryByIDFunc(id)
	}
	return &model.Category{}, nil
}

func (m *MockCategoryRepository) UpdateCategory(category *model.Category) error {
	if m.UpdateCategoryFunc != nil {
		return m.UpdateCategoryFunc(category)
	}
	return nil
}

func (m *MockCategoryRepository) DeleteCategory(id string) error {
	if m.DeleteCategoryFunc != nil {
		return m.DeleteCategoryFunc(id)
	}
	return nil
}

func TestCategoryService_CreateCategory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockCategoryRepository{
			CreateCategoryFunc: func(c *model.Category) error { return nil },
		}
		service := NewCategoryService(mockRepo)

		cat := &model.Category{Name: "Tech"}
		err := service.CreateCategory(cat)

		if err != nil {
			t.Errorf("Se esperaba error nil, se obtuvo %v", err)
		}
	})

	t.Run("NilBody", func(t *testing.T) {
		service := NewCategoryService(&MockCategoryRepository{})
		err := service.CreateCategory(nil)

		if err == nil || err.Error() != "El Body Request viene vacío." {
			t.Errorf("Se esperaba error 'El Body Request viene vacío.', se obtuvo %v", err)
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo := &MockCategoryRepository{
			CreateCategoryFunc: func(c *model.Category) error { return errors.New("DB error") },
		}
		service := NewCategoryService(mockRepo)

		err := service.CreateCategory(&model.Category{})
		if err == nil || err.Error() != "Error al crear la categoria" {
			t.Errorf("Se esperaba error de servicio, se obtuvo %v", err)
		}
	})
}

func TestCategoryService_FindAllCategories(t *testing.T) {
	t.Run("EmptyList", func(t *testing.T) {
		mockRepo := &MockCategoryRepository{
			FindAllCategoriesFunc: func() ([]model.Category, error) {
				return []model.Category{}, nil
			},
		}
		service := NewCategoryService(mockRepo)

		_, err := service.FindAllCategories()
		if err == nil || err.Error() != "No se encontraron categorias" {
			t.Errorf("Se esperaba error 'No se encontraron categorias', se obtuvo %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockCategoryRepository{
			FindAllCategoriesFunc: func() ([]model.Category, error) {
				return []model.Category{{Name: "Tech"}}, nil
			},
		}
		service := NewCategoryService(mockRepo)

		cats, err := service.FindAllCategories()
		if err != nil {
			t.Errorf("No se esperaba error, se obtuvo %v", err)
		}
		if len(cats) != 1 {
			t.Errorf("Se esperaba 1 categoria, se obtuvo %d", len(cats))
		}
	})
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	t.Run("ErrorFinding", func(t *testing.T) {
		mockRepo := &MockCategoryRepository{
			FindCategoryByIDFunc: func(id string) (*model.Category, error) {
				return nil, errors.New("not found")
			},
		}
		service := NewCategoryService(mockRepo)

		err := service.DeleteCategory("123")
		if err == nil || err.Error() != "Error al eliminar la categoría" {
			t.Errorf("Se esperaba error al eliminar, se obtuvo %v", err)
		}
	})
}
