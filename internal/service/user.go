package service

import (
	"errors"

	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/repository"
	"github.com/lucas-dev/in-memory-db/pkg/utils"
)

type UserService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) RegisterUser(req *model.UserCreateRequest) error {
	if req == nil {
		return errors.New("el cuerpo de la solicitud está vacío")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("Error al hashear la contraseña: " + err.Error())
	}

	user := model.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := s.userRepo.RegisterUser(&user); err != nil {
		return errors.New("Error al crear el usuario: " + err.Error())
	}

	return nil
}

func (s *UserService) LoginUser(req *model.LoginUserRequest) (string, error) {
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	err = utils.ComparePassword(user.Password, req.Password)
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	token, err := utils.GenerateToken(s.jwtSecret, user.ID.String(), user.Email, user.FirstName, user.LastName)
	if err != nil {
		return "", errors.New("error al generar el token")
	}

	return token, nil
}

func (s *UserService) FindUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindUserByID(id)
	if err != nil {
		return nil, errors.New("Error al buscar el usuario: " + err.Error())
	}
	return user, nil
}

func (s *UserService) UpdateUser(user *model.User) error {
	err := s.userRepo.UpdateUser(user)
	if err != nil {
		return errors.New("Error al actualizar el usuario: " + err.Error())
	}

	return nil
}

func (s *UserService) DeleteUser(email string) error {
	err := s.userRepo.DeleteUser(email)
	if err != nil {
		return errors.New("Error al eliminar el usuario: " + err.Error())
	}

	return nil
}
