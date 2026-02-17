package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/service"
	"github.com/lucas-dev/in-memory-db/pkg/utils"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder godoc
// @Summary Crear orden
// @Description Crear orden de pago
// @Tags orders
// @Produce json
// @Param body body model.Order true "Datos de la orden"
// @Success 201 {object} model.Order
// @Failure 400 {object} map[string]string "JSON invalido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, ok := h.getAuthUserID(c)
	if !ok {
		return
	}

	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalido: " + err.Error()})
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

// GetAllOrders godoc
// @Summary Obtener todas las ordenes
// @Description Obtener todas las ordenes de pago
// @Tags orders
// @Produce json
// @Success 200 {object} []model.Order
// @Failure 500 {object} map[string]string "Error interno"
// @Router /orders [get]
func (h *OrderHandler) GetAllOrders(g *gin.Context) {
	authUserID, ok := h.getAuthUserID(g)
	if !ok {
		return
	}

	orders, err := h.orderService.GetOrdersByUserID(authUserID.String())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, orders)
}

// GetOrderByID godoc
// @Summary Obtener order
// @Description Obtener orden de pago por UUID
// @Tags orders
// @Produce json
// @Param id path string true "UUID de la orden"
// @Success 200 {object} model.Order
// @Failure 400 {object} map[string]string "UUID invalido"
// @Failure 500 {object} map[string]string "Error interno"
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByID(g *gin.Context) {
	orderID := g.Param("id")

	if orderID == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	if _, err := uuid.Parse(orderID); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalido"})
		return
	}

	order, ok := h.getOwnedOrder(g, orderID, "No tenes permisos para ver esta orden")
	if !ok {
		return
	}

	g.JSON(http.StatusOK, order)
}

// func (h *OrderHandler) DeleteOrder(c *gin.Context) {
// 	orderID := c.Param("id")
//
// 	if _, err := uuid.Parse(orderID); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID invalido"})
// 		return
// 	}
//
// 	err := h.orderService.DeleteOrder(orderID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{"message": "Orden eliminada exitosamente"})
// }

// UpdateStatusOrder godoc
// @Summary Actualizar estado de orden
// @Description Actualizar estado de pago de una orden segun su UUID
// @Tags orders
// @Produce json
// @Param id path string true "UUID de la orden"
// @Param body body model.UpdateOrderStatusRequest true "Datos de la orden"
// @Success 200 {object} map[string]string "Orden actualizada exitosamente"
// @Failure 400 {object} map[string]string "UUID invalido o error al actualizar"
// @Router /orders/{id} [patch]
func (h *OrderHandler) UpdateStatusOrder(c *gin.Context) {
	orderID := c.Param("id")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, "Debe mandar un OrderID")
		return
	}

	_, err := uuid.Parse(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, "UUID invalido")
		return
	}

	if _, ok := h.getOwnedOrder(c, orderID, "No tenes permisos para modificar esta orden"); !ok {
		return
	}

	var newOrder model.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, "Solo puedo modificar propiedades validas de la orden")
		return
	}

	if !utils.IsValidOrderStatus(newOrder.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Estado invalido. Debe ser: PENDING, PAID, SHIPPED o CANCELLED"})
		return
	}

	err = h.orderService.UpdateStatusOrder(&newOrder, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Orden actualizada exitosamente"})
}

// OrderUpdatePay godoc
// @Summary Cambiar estado de orden a pagada
// @Description Cambiar estao de la orden a pagada segun su UUID
// @Tags orders
// @Produce json
// @Param id path string true "UUID de la orden"
// @Success 200 {object} map[string]string "Orden pagada exitosamente"
// @Failure 400 {object} map[string]string "UUID invalido o error al pagar"
// @Router /orders/{id}/pay [patch]
func (h *OrderHandler) OrderUpdatePay(c *gin.Context) {
	orderID := c.Param("id")

	if _, err := uuid.Parse(orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID de la orden no valido"})
		return
	}

	if _, ok := h.getOwnedOrder(c, orderID, "No tenes permisos para pagar esta orden"); !ok {
		return
	}

	if errOrderPay := h.orderService.OrderUpdatePay(orderID); errOrderPay != nil {
		switch {
		case errors.Is(errOrderPay, service.ErrInsufficientStock):
			c.JSON(http.StatusConflict, gin.H{"error": "Stock insuficiente para pagar la orden"})
		case errors.Is(errOrderPay, service.ErrOrderNotPending):
			c.JSON(http.StatusConflict, gin.H{"error": "La orden no esta en estado PENDING"})
		case errors.Is(errOrderPay, service.ErrOrderNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "No se encontro la orden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo cambiar el estado de la orden a pagada"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Orden pagada exitosamente"})
}

func (h *OrderHandler) getAuthUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return uuid.Nil, false
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario invalido en token"})
		return uuid.Nil, false
	}

	authUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID de usuario invalido"})
		return uuid.Nil, false
	}

	return authUserID, true
}

func (h *OrderHandler) getOwnedOrder(c *gin.Context, orderID, forbiddenMessage string) (*model.Order, bool) {
	authUserID, ok := h.getAuthUserID(c)
	if !ok {
		return nil, false
	}

	order, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontro la orden"})
		return nil, false
	}

	if order.UserID != authUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": forbiddenMessage})
		return nil, false
	}

	return order, true
}
