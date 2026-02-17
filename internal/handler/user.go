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

// RegisterUser godoc
// @Summary Registrar usuario
// @Description Crea un nuevo usuario
// @Tags users
// @Accept json
// @Produce json
// @Param body body model.UserCreateRequest true "Datos de registro"
// @Success 201 {object} map[string]string "Usuario creado exitosamente"
// @Failure 400 {object} map[string]string "JSON invalido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /users/register [post]
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

// LoginUser godoc
// @Summary Login de usuario
// @Description Autentica usuario y devuelve JWT
// @Tags users
// @Accept json
// @Produce json
// @Param body body model.LoginUserRequest true "Credenciales"
// @Success 200 {object} map[string]string "message + token"
// @Failure 400 {object} map[string]string "JSON invalido"
// @Failure 401 {object} map[string]string "Credenciales invalidas"
// @Router /users/login [post]
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

// UpdateUser godoc
// @Summary Actualizar usuario
// @Description Actualiza datos del usuario
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.User true "Usuario actualizado"
// @Success 200 {object} map[string]interface{} "Usuario actualizado"
// @Failure 400 {object} map[string]string "JSON invalido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /users [patch]
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

// DeleteUser godoc
// @Summary Eliminar usuario
// @Description Elimina un usuario por email
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param email path string true "Email del usuario"
// @Success 200 {object} map[string]string "Usuario eliminado"
// @Failure 400 {object} map[string]string "Email requerido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /users/{email} [delete]
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
