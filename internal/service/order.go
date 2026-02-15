package service

import (
	"errors"

	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/repository"
)

type OrderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

func (s *OrderService) CreateOrder(order *model.Order) (*model.Order, error) {
	if order == nil {
		return nil, errors.New("El Body Request viene vacío.")
	}

	order, err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetAllOrders() ([]model.Order, error) {
	res, err := s.orderRepo.GetAllOrders()
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("No se encontraron órdenes")
	}

	return res, nil
}

func (s *OrderService) GetOrderByID(id string) (*model.Order, error) {
	res, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return nil, errors.New("Error al obtener la orden")
	}

	if res == nil {
		return nil, errors.New("No se encontró la orden")
	}

	return res, nil

}

func (s *OrderService) UpdateOrder(order *model.Order, id string) (*model.Order, error) {
	_, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return nil, errors.New("Error al obtener la orden")
	}

	return s.orderRepo.UpdateOrder(order, id)
}

func (s *OrderService) DeleteOrder(id string) error {
	_, err := s.orderRepo.GetOrderByID(id)
	if err != nil {
		return errors.New("Error al obtener la orden")
	}

	return s.orderRepo.DeleteOrder(id)
}
