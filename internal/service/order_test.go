package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/repository"
	"gorm.io/gorm"
)

type MockOrderRepository struct {
	CreateOrderFunc       func(order *model.Order) (*model.Order, error)
	GetAllOrdersFunc      func() ([]model.Order, error)
	GetOrdersByUserIDFunc func(userID string) ([]model.Order, error)
	GetOrderByIDFunc      func(id string) (*model.Order, error)
	UpdateStatusOrderFunc func(req *model.UpdateOrderStatusRequest, id string) error
	OrderUpdatePayFunc    func(id string) error
}

func (m *MockOrderRepository) CreateOrder(order *model.Order) (*model.Order, error) {
	if m.CreateOrderFunc != nil {
		return m.CreateOrderFunc(order)
	}
	return order, nil
}

func (m *MockOrderRepository) GetAllOrders() ([]model.Order, error) {
	if m.GetAllOrdersFunc != nil {
		return m.GetAllOrdersFunc()
	}
	return []model.Order{}, nil
}

func (m *MockOrderRepository) GetOrdersByUserID(userID string) ([]model.Order, error) {
	if m.GetOrdersByUserIDFunc != nil {
		return m.GetOrdersByUserIDFunc(userID)
	}
	return []model.Order{}, nil
}

func (m *MockOrderRepository) GetOrderByID(id string) (*model.Order, error) {
	if m.GetOrderByIDFunc != nil {
		return m.GetOrderByIDFunc(id)
	}
	return &model.Order{}, nil
}

func (m *MockOrderRepository) UpdateStatusOrder(req *model.UpdateOrderStatusRequest, id string) error {
	if m.UpdateStatusOrderFunc != nil {
		return m.UpdateStatusOrderFunc(req, id)
	}
	return nil
}

func (m *MockOrderRepository) OrderUpdatePay(id string) error {
	if m.OrderUpdatePayFunc != nil {
		return m.OrderUpdatePayFunc(id)
	}
	return nil
}

func TestOrderService_CreateOrder(t *testing.T) {
	t.Run("NilBody", func(t *testing.T) {
		service := NewOrderService(&MockOrderRepository{})
		_, err := service.CreateOrder(nil)
		if err == nil || err.Error() != "El Body Request viene vacío." {
			t.Errorf("Se esperaba error body vacío, se obtuvo %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		expectedID := uuid.New()

		mockRepo := &MockOrderRepository{
			CreateOrderFunc: func(o *model.Order) (*model.Order, error) {
				o.ID = expectedID
				return o, nil
			},
		}
		service := NewOrderService(mockRepo)

		order := &model.Order{}

		res, err := service.CreateOrder(order)
		if err != nil {
			t.Errorf("No se esperaba error, se obtuvo %v", err)
		}

		if res.ID != expectedID {
			t.Errorf("Se esperaba ID %v, se obtuvo %v", expectedID, res.ID)
		}
	})
}

func TestOrderService_GetOrderByID(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		mockRepo := &MockOrderRepository{
			GetOrderByIDFunc: func(id string) (*model.Order, error) {
				return nil, nil
			},
		}
		service := NewOrderService(mockRepo)

		_, err := service.GetOrderByID("123")
		if err == nil || err.Error() != "No se encontró la orden" {
			t.Errorf("Se esperaba error 'No se encontró la orden', se obtuvo %v", err)
		}
	})

	t.Run("DBError", func(t *testing.T) {
		mockRepo := &MockOrderRepository{
			GetOrderByIDFunc: func(id string) (*model.Order, error) {
				return nil, errors.New("connection lost")
			},
		}
		service := NewOrderService(mockRepo)

		_, err := service.GetOrderByID("123")
		if err == nil || err.Error() != "Error al obtener la orden" {
			t.Errorf("Se esperaba error genérico, se obtuvo %v", err)
		}
	})
}

func TestOrderService_OrderUpdatePay(t *testing.T) {
	tests := []struct {
		name      string
		repoError error
		svcError  error
	}{
		{
			name:      "Success",
			repoError: nil,
			svcError:  nil,
		},
		{
			name:      "Insufficient Stock",
			repoError: repository.ErrInsufficientStock,
			svcError:  ErrInsufficientStock,
		},
		{
			name:      "Order Not Pending",
			repoError: repository.ErrOrderNotPending,
			svcError:  ErrOrderNotPending,
		},
		{
			name:      "Record Not Found (GORM)",
			repoError: gorm.ErrRecordNotFound,
			svcError:  ErrOrderNotFound,
		},
		{
			name:      "Unknown Error",
			repoError: errors.New("db crash"),
			svcError:  errors.New("Error al actualizar el estado de la orden a pagada"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockOrderRepository{
				OrderUpdatePayFunc: func(id string) error {
					return tt.repoError
				},
			}
			service := NewOrderService(mockRepo)

			err := service.OrderUpdatePay("123")

			if tt.svcError == nil && err != nil {
				t.Errorf("Se esperaba nil, se obtuvo %v", err)
			}
			if tt.svcError != nil && err == nil {
				t.Errorf("Se esperaba error %v, se obtuvo nil", tt.svcError)
			}
			if tt.svcError != nil && err.Error() != tt.svcError.Error() {
				t.Errorf("Se esperaba error %v, se obtuvo %v", tt.svcError, err)
			}
		})
	}
}
