package service

import (
	"errors"

	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/repository"
)

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) CreateProduct(product *model.Product) error {
	if product == nil {
		return errors.New("El Body Request viene vacío.")
	}

	return s.productRepo.CreateProduct(product)
}

func (s *ProductService) GetAllProducts(limit, page int) (*model.ProductResponse, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	_, total, err := s.productRepo.GetAllProducts(0, 0)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	offset := (page - 1) * limit

	prods, total, err := s.productRepo.GetAllProducts(limit, offset)
	if err != nil {
		return nil, err
	}

	response := &model.ProductResponse{
		Data:       prods,
		Total:      int(total),
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	return response, nil
}

func (s *ProductService) GetProductByID(id string) (model.Product, error) {
	prod, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return model.Product{}, err
	}

	if prod == nil {
		return model.Product{}, errors.New("No se encontró el producto")
	}
	return *prod, nil
}

func (s *ProductService) UpdateProduct(product *model.Product, id string) error {
	if product == nil {
		return errors.New("El Body Request viene vacío.")
	}

	_, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return errors.New("Error al obtener el producto: " + err.Error())
	}

	err = s.productRepo.UpdateProduct(product, id)
	if err != nil {
		return errors.New("Error al actualizar el producto: " + err.Error())
	}

	return nil
}

func (s *ProductService) DeleteProduct(id string) error {
	if id == "" {
		return errors.New("El ID viene vacío.")
	}

	_, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return errors.New("Error al obtener el producto: " + err.Error())
	}

	err = s.productRepo.DeleteProduct(id)
	if err != nil {
		return errors.New("Error al eliminar el producto: " + err.Error())
	}
	return nil
}
