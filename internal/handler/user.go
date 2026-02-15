package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req model.UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalido: " + err.Error()})
		return
	}

	if err := h.userService.RegisterUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado exitosamente"})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req model.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalido: " + err.Error()})
		return
	}

	token, err := h.userService.LoginUser(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sesion iniciada exitosamente",
		"token":   token,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalido: " + err.Error()})
		return
	}

	if err := h.userService.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado exitosamente", "user": user})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	email := c.Param("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email es requerido"})
		return
	}

	if err := h.userService.DeleteUser(email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}
