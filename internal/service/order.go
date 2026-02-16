package service

import (
	"errors"

	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/repository"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo repository.OrderRepository
}

var (
	ErrOrderNotFound     = errors.New("orden no encontrada")
	ErrOrderNotPending   = errors.New("la orden no está en estado PENDING")
	ErrInsufficientStock = errors.New("stock insuficiente para completar el pago")
)

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

func (s *OrderService) GetOrdersByUserID(userID string) ([]model.Order, error) {
	if userID == "" {
		return nil, errors.New("ID de usuario vacío")
	}

	res, err := s.orderRepo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return []model.Order{}, nil
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

// func (s *OrderService) DeleteOrder(id string) error {
// 	_, err := s.orderRepo.GetOrderByID(id)
// 	if err != nil {
// 		return errors.New("Error al obtener la orden")
// 	}

// 	return s.orderRepo.DeleteOrder(id)
// }

func (s *OrderService) UpdateStatusOrder(order *model.UpdateOrderStatusRequest, id string) error {
	if _, err := s.orderRepo.GetOrderByID(id); err != nil {
		return errors.New("Error al obtener la orden")
	}

	if err := s.orderRepo.UpdateStatusOrder(&model.UpdateOrderStatusRequest{Status: order.Status}, id); err != nil {
		return errors.New("Error al actualizar el estado de la orden")
	}

	return nil
}

func (s *OrderService) OrderUpdatePay(id string) error {
	if err := s.orderRepo.OrderUpdatePay(id); err != nil {
		switch {
		case errors.Is(err, repository.ErrInsufficientStock):
			return ErrInsufficientStock
		case errors.Is(err, repository.ErrOrderNotPending):
			return ErrOrderNotPending
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrOrderNotFound
		default:
			return errors.New("Error al actualizar el estado de la orden a pagada")
		}
	}

	return nil
}
