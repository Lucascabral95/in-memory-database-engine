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
