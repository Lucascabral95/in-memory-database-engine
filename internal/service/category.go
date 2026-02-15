package service

import (
	"errors"

	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/repository"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(category *model.Category) error {
	if category == nil {
		return errors.New("El Body Request viene vacío.")
	}

	err := s.categoryRepo.CreateCategory(category)
	if err != nil {
		return errors.New("Error al crear la categoria")
	}

	return nil
}

func (s *CategoryService) FindAllCategories() ([]model.Category, error) {
	categories, err := s.categoryRepo.FindAllCategories()
	if err != nil {
		return nil, errors.New("Error al obtener las categorias")
	}

	if len(categories) == 0 {
		return nil, errors.New("No se encontraron categorias")
	}

	return categories, nil
}

func (s *CategoryService) FindCategoryByID(id string) (*model.Category, error) {
	category, err := s.categoryRepo.FindCategoryByID(id)
	if err != nil {
		return nil, errors.New("Error al obtener la categoría")
	}

	return category, nil
}

func (s *CategoryService) UpdateCategory(category *model.Category) error {
	if category == nil {
		return errors.New("El Body Request viene vacío.")
	}

	err := s.categoryRepo.UpdateCategory(category)
	if err != nil {
		return errors.New("Error al actualizar la categoría")
	}

	return nil
}

func (s *CategoryService) DeleteCategory(id string) error {
	_, errFindById := s.categoryRepo.FindCategoryByID(id)
	if errFindById != nil {
		return errors.New("Error al eliminar la categoría")
	}

	err := s.categoryRepo.DeleteCategory(id)
	if err != nil {
		return errors.New("Error al eliminar la categoría")
	}

	return nil
}
