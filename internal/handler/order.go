package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido en token"})
		return
	}

	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	var totalAmountForProducts float64

	for _, price := range order.Items {
		totalAmountForProducts += price.PriceAtMoment * float64(price.Quantity)
	}

	order.UserID = userID
	order.Status = model.OrderStatusPending
	order.TotalAmount = totalAmountForProducts

	createdOrder, err := h.orderService.CreateOrder(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
}

func (h *OrderHandler) GetAllOrders(g *gin.Context) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderByID(g *gin.Context) {
	orderID := g.Param("id")

	if orderID == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if _, err := uuid.Parse(orderID); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "UUID inválido"})
		return
	}

	order, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	orderID := c.Param("id")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, "Debe mandar un OrderID")
	}

	_, err := uuid.Parse(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "UUID inválido")
		return
	}

	var newOrder model.Order
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, "Solo puedo modificar propiedades válidas de la orden")
		return
	}

	res, err := h.orderService.UpdateOrder(&newOrder, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Orden actualizada exitosamente", "data": res})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")

	if _, err := uuid.Parse(orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID inválido"})
		return
	}

	err := h.orderService.DeleteOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Orden eliminada exitosamente"})
}
