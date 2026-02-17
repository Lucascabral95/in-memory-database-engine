package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/service"
	"github.com/lucas-dev/in-memory-db/pkg/utils"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct godoc
// @Summary Crear producto
// @Description Crea un producto nuevo y genera SKU automaticamente.
// @Tags products
// @Accept json
// @Produce json
// @Param body body model.Product true "Datos del producto"
// @Success 201 {object} model.Product
// @Failure 400 {object} map[string]string "JSON invalido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product model.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	product.SKU = utils.GenerateSKU()

	if err := h.productService.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetAllProducts godoc
// @Summary Listar productos
// @Description Lista productos paginados.
// @Tags products
// @Produce json
// @Param limit query int false "Cantidad por pagina" default(10)
// @Param page query int false "Numero de pagina" default(1)
// @Success 200 {object} model.ProductResponse
// @Failure 500 {object} map[string]string "Error interno"
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := utils.ConvertAnyToInt(limitStr)

	page, err := utils.ConvertAnyToInt(pageStr)

	products, err := h.productService.GetAllProducts(limit, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        products.Data,
		"total":       products.Total,
		"page":        products.Page,
		"limit":       products.Limit,
		"total_pages": products.TotalPages,
	})
}

// GetProductByID godoc
// @Summary Obtener producto
// @Description Devuelve un producto por su UUID.
// @Tags products
// @Produce json
// @Param id path string true "UUID del producto"
// @Success 200 {object} model.Product
// @Failure 400 {object} map[string]string "UUID invalido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID := c.Param("id")

	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if _, err := uuid.Parse(productID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID inválido"})
		return
	}

	product, err := h.productService.GetProductByID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Actualizar producto
// @Description Actualiza nombre, descripcion, precio, stock o categoria de un producto.
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "UUID del producto"
// @Param body body model.Product true "Campos a actualizar"
// @Success 200 {object} map[string]string "Producto actualizado exitosamente"
// @Failure 400 {object} map[string]string "UUID/JSON invalido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /products/{id} [patch]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")

	if _, err := uuid.Parse(productID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID inválido"})
		return
	}

	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	err := h.productService.UpdateProduct(&product, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto actualizado exitosamente"})
}

// DeleteProduct godoc
// @Summary Eliminar producto
// @Description Elimina un producto por UUID.
// @Tags products
// @Produce json
// @Param id path string true "UUID del producto"
// @Success 200 {object} map[string]string "Producto eliminado exitosamente"
// @Failure 400 {object} map[string]string "UUID invalido o error al eliminar"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productId := c.Param("id")

	if _, err := uuid.Parse(productId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID inválido"})
		return
	}

	err := h.productService.DeleteProduct(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al eliminar el producto"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado exitosamente"})
}
