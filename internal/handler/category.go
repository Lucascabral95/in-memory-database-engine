package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/service"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Crear categoria
// @Description Crear una categoria nueva de productos
// @Tags categories
// @Accept json
// @Produce json
// @Param body body model.Category true "Datos de la categoria"
// @Success 201 {object} model.Category
// @Failure 400 {object} map[string]string "JSON invalido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category model.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	if err := h.categoryService.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// FindAllCategories godoc
// @Summary Buscar categorias
// @Description Buscar todas las categorias de la API
// @Tags categories
// @Produce json
// @Success 200 {object} []model.Category
// @Failure 500 {object} map[string]string "Error interno"
// @Router /categories [get]
func (h *CategoryHandler) FindAllCategories(c *gin.Context) {
	categories, err := h.categoryService.FindAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// FindCategoryByID godoc
// @Summary Buscar categoria
// @Description Buscar una categoria por su UUID
// @Tags categories
// @Produce json
// @Param id path string true "UUID de la categoria"
// @Success 200 {object} model.Category
// @Failure 400 {object} map[string]string "UUID invalido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /categories/{id} [get]
func (h *CategoryHandler) FindCategoryByID(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID inválido"})
		return
	}

	category, err := h.categoryService.FindCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// UpdateCategory godoc
// @Summary Actualizar categoría
// @Description Actualizar nombre de la categoria segun su UUID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "UUID de la categoría"
// @Param body body model.Category true "Datos de la categoría"
// @Success 200 {object} model.Category
// @Failure 400 {object} map[string]string "UUID/JSON invalido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /categories/{id} [patch]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	categoryID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID inválido"})
		return
	}

	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	category.ID = categoryID

	if err := h.categoryService.UpdateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary Eliminar producto
// @Description Eliminar un producto segun su UUID
// @Tags categories
// @Produce json
// @Param id path string true "UUID de la categoría"
// @Success 200 {object} map[string]string "Categoría eliminada exitosamente"
// @Failure 400 {object} map[string]string "UUID invalido o error al eliminar"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID inválido"})
		return
	}

	if err := h.categoryService.DeleteCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Categoría eliminada exitosamente"})
}

func (h *CategoryHandler) RegisterRoutes(router *gin.Engine) {
	categories := router.Group("/categories")
	{
		categories.POST("", h.CreateCategory)
		categories.GET("", h.FindAllCategories)
		categories.GET("/:id", h.FindCategoryByID)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}
}
