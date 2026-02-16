package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/pkg/utils"
)

type MockUserRepository struct {
	RegisterUserFunc    func(user *model.User) error
	FindUserByEmailFunc func(email string) (*model.User, error)
	FindUserByIDFunc    func(id uint) (*model.User, error)
	UpdateUserFunc      func(user *model.User) error
	DeleteUserFunc      func(email string) error
}

func (m *MockUserRepository) RegisterUser(user *model.User) error {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(user)
	}
	return nil
}

func (m *MockUserRepository) FindUserByEmail(email string) (*model.User, error) {
	if m.FindUserByEmailFunc != nil {
		return m.FindUserByEmailFunc(email)
	}
	return &model.User{}, nil
}

func (m *MockUserRepository) FindUserByID(id uint) (*model.User, error) {
	if m.FindUserByIDFunc != nil {
		return m.FindUserByIDFunc(id)
	}
	return &model.User{}, nil
}

func (m *MockUserRepository) UpdateUser(user *model.User) error {
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(user)
	}
	return nil
}

func (m *MockUserRepository) DeleteUser(email string) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(email)
	}
	return nil
}

func TestUserService_RegisterUser(t *testing.T) {
	t.Run("NilBody", func(t *testing.T) {
		svc := NewUserService(&MockUserRepository{}, "secret")
		err := svc.RegisterUser(nil)
		if err == nil || err.Error() != "el cuerpo de la solicitud está vacío" {
			t.Errorf("Se esperaba error body vacío, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			RegisterUserFunc: func(u *model.User) error {
				if u.Password == "password123" {
					t.Error("La contraseña no fue hasheada")
				}
				return nil
			},
		}
		svc := NewUserService(mockRepo, "secret")
		req := &model.UserCreateRequest{
			Email:     "test@test.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}
		err := svc.RegisterUser(req)
		if err != nil {
			t.Errorf("No se esperaba error, got %v", err)
		}
	})
}

func TestUserService_LoginUser(t *testing.T) {
	hashedPass, _ := utils.HashPassword("mypassword")
	mockUUID := uuid.New()

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			FindUserByEmailFunc: func(email string) (*model.User, error) {
				return nil, errors.New("not found")
			},
		}
		svc := NewUserService(mockRepo, "secret")
		_, err := svc.LoginUser(&model.LoginUserRequest{Email: "fail@fail.com"})
		if err == nil || err.Error() != "credenciales inválidas" {
			t.Errorf("Se esperaba error credenciales inválidas, got %v", err)
		}
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			FindUserByEmailFunc: func(email string) (*model.User, error) {
				return &model.User{Password: hashedPass}, nil
			},
		}
		svc := NewUserService(mockRepo, "secret")
		_, err := svc.LoginUser(&model.LoginUserRequest{Email: "test@test.com", Password: "wrongpass"})
		if err == nil || err.Error() != "credenciales inválidas" {
			t.Errorf("Se esperaba error credenciales inválidas, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			FindUserByEmailFunc: func(email string) (*model.User, error) {
				return &model.User{
					BaseModel: model.BaseModel{
						ID: mockUUID,
					},
					Email:     "test@test.com",
					Password:  hashedPass,
					FirstName: "Test",
					LastName:  "User",
				}, nil
			},
		}
		svc := NewUserService(mockRepo, "super-secret-key")
		token, err := svc.LoginUser(&model.LoginUserRequest{Email: "test@test.com", Password: "mypassword"})

		if err != nil {
			t.Fatalf("No se esperaba error, got %v", err)
		}
		if token == "" {
			t.Error("Se esperaba un token JWT")
		}

		claims, _ := utils.ValidateToken("super-secret-key", token)
		if claims.Email != "test@test.com" {
			t.Error("El token no contiene el email correcto")
		}
	})
}

func TestUserService_FindUserByID(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			FindUserByIDFunc: func(id uint) (*model.User, error) {
				return nil, errors.New("db error")
			},
		}
		svc := NewUserService(mockRepo, "secret")
		_, err := svc.FindUserByID(1)
		if err == nil {
			t.Error("Se esperaba error")
		}
	})
}
