package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucas-dev/in-memory-db/internal/model"
	"github.com/lucas-dev/in-memory-db/internal/service"
	"github.com/lucas-dev/in-memory-db/internal/storage"
)

const cartTTLSeconds = 24 * 60 * 60

type CartHandler struct {
	store          *storage.MemoryStore
	productService *service.ProductService
	orderService   *service.OrderService
}

func NewCartHandler(
	store *storage.MemoryStore,
	productService *service.ProductService,
	orderService *service.OrderService,
) *CartHandler {
	return &CartHandler{
		store:          store,
		productService: productService,
		orderService:   orderService,
	}
}

type addToCartRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type updateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, userIDStr, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req addToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalido: " + err.Error()})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity debe ser mayor a 0"})
		return
	}

	prodID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID de producto invalido"})
		return
	}

	product, err := h.productService.GetProductByID(prodID.String())
	if err != nil {
		status := http.StatusInternalServerError
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": "No se pudo obtener el producto: " + err.Error()})
		return
	}

	key := cartKey(userIDStr)
	cart, exists, err := h.loadCart(key, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo decodificar el carrito en memoria"})
		return
	}

	targetQuantity := req.Quantity
	for i := range cart.Items {
		if cart.Items[i].ProductID == prodID {
			targetQuantity += cart.Items[i].Quantity
			break
		}
	}

	if targetQuantity > product.Stock {
		c.JSON(http.StatusConflict, gin.H{
			"error":              "Stock insuficiente para agregar al carrito",
			"product_id":         prodID.String(),
			"requested_quantity": targetQuantity,
			"available_stock":    product.Stock,
		})
		return
	}

	merged := false
	for i := range cart.Items {
		if cart.Items[i].ProductID == prodID {
			cart.Items[i].Quantity += req.Quantity
			cart.Items[i].Price = product.Price
			merged = true
			break
		}
	}

	if !merged {
		cart.Items = append(cart.Items, model.CartItem{
			ProductID: prodID,
			Quantity:  req.Quantity,
			Price:     product.Price,
		})
	}

	if err := h.saveCart(key, cart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo serializar el carrito"})
		return
	}

	status := "saved in RAM"
	if exists {
		status = "updated in RAM"
	}

	c.JSON(http.StatusOK, gin.H{"status": status, "cart": cart})
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, userIDStr, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID de producto invalido"})
		return
	}

	var req updateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalido: " + err.Error()})
		return
	}

	key := cartKey(userIDStr)
	cart, exists, err := h.loadCart(key, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo decodificar el carrito en memoria"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart expired or empty"})
		return
	}

	itemIndex := -1
	for i := range cart.Items {
		if cart.Items[i].ProductID == productID {
			itemIndex = i
			break
		}
	}

	if itemIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "El producto no existe en el carrito"})
		return
	}

	if req.Quantity <= 0 {
		cart.Items = append(cart.Items[:itemIndex], cart.Items[itemIndex+1:]...)
		if len(cart.Items) == 0 {
			h.store.Del(key)
		} else {
			if err := h.saveCart(key, cart); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo serializar el carrito"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "item_removed",
			"cart":   cart,
		})
		return
	}

	product, err := h.productService.GetProductByID(productID.String())
	if err != nil {
		status := http.StatusInternalServerError
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": "No se pudo obtener el producto: " + err.Error()})
		return
	}

	if req.Quantity > product.Stock {
		c.JSON(http.StatusConflict, gin.H{
			"error":              "Stock insuficiente para actualizar el item",
			"product_id":         productID.String(),
			"requested_quantity": req.Quantity,
			"available_stock":    product.Stock,
		})
		return
	}

	cart.Items[itemIndex].Quantity = req.Quantity
	cart.Items[itemIndex].Price = product.Price

	if err := h.saveCart(key, cart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo serializar el carrito"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "item_updated",
		"cart":   cart,
	})
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	userID, userIDStr, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID de producto invalido"})
		return
	}

	key := cartKey(userIDStr)
	cart, exists, err := h.loadCart(key, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo decodificar el carrito en memoria"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart expired or empty"})
		return
	}

	removed := false
	filtered := make([]model.CartItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		if item.ProductID == productID {
			removed = true
			continue
		}
		filtered = append(filtered, item)
	}

	if !removed {
		c.JSON(http.StatusNotFound, gin.H{"error": "El producto no existe en el carrito"})
		return
	}

	cart.Items = filtered
	if len(cart.Items) == 0 {
		h.store.Del(key)
	} else {
		if err := h.saveCart(key, cart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo serializar el carrito"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "item_removed",
		"cart":   cart,
	})
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, userIDStr, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	h.store.Del(cartKey(userIDStr))
	c.JSON(http.StatusOK, gin.H{
		"status": "cart_cleared",
		"cart": model.Cart{
			UserID: userID,
			Items:  []model.CartItem{},
		},
	})
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, userIDStr, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	key := cartKey(userIDStr)
	cart, exists, err := h.loadCart(key, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo decodificar el carrito en memoria"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart expired or empty"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) Checkout(c *gin.Context) {
	userID, userIDStr, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	key := cartKey(userIDStr)
	cart, exists, err := h.loadCart(key, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, checkoutResponse("failed", nil, "", "", "No se pudo decodificar el carrito en memoria"))
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, checkoutResponse("failed", nil, "", "", "Cart expired or empty"))
		return
	}

	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, checkoutResponse("failed", nil, "", "", "No hay productos en el carrito"))
		return
	}

	orderItems := make([]model.OrderItem, 0, len(cart.Items))
	totalAmount := 0.0

	for _, item := range cart.Items {
		if item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, checkoutResponse("failed", nil, "", "", "El carrito tiene cantidades invalidas"))
			return
		}

		product, err := h.productService.GetProductByID(item.ProductID.String())
		if err != nil {
			status := http.StatusInternalServerError
			if isNotFoundError(err) {
				status = http.StatusNotFound
			}

			c.JSON(status, checkoutResponse("failed", nil, "", "", "No se pudo obtener el producto para checkout: "+err.Error()))
			return
		}

		orderItems = append(orderItems, model.OrderItem{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			PriceAtMoment: product.Price,
		})

		totalAmount += product.Price * float64(item.Quantity)
	}

	order := &model.Order{
		UserID:      userID,
		Status:      model.OrderStatusPending,
		TotalAmount: totalAmount,
		Items:       orderItems,
	}

	createdOrder, err := h.orderService.CreateOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, checkoutResponse("failed", nil, "", "", "No se pudo crear la orden: "+err.Error()))
		return
	}

	orderID := createdOrder.ID

	if err := h.orderService.OrderUpdatePay(orderID.String()); err != nil {
		switch {
		case errors.Is(err, service.ErrInsufficientStock):
			orderStatus := model.OrderStatusPending
			warning := ""

			cancelErr := h.orderService.UpdateStatusOrder(
				&model.UpdateOrderStatusRequest{Status: model.OrderStatusCancelled},
				orderID.String(),
			)
			if cancelErr == nil {
				orderStatus = model.OrderStatusCancelled
			} else {
				warning = "No se pudo cancelar la orden automaticamente"
			}

			resp := checkoutResponse("failed", &orderID, orderStatus, "", "Stock insuficiente para completar el checkout")
			if warning != "" {
				resp["warning"] = warning
			}
			c.JSON(http.StatusConflict, resp)
		case errors.Is(err, service.ErrOrderNotPending):
			c.JSON(http.StatusConflict, checkoutResponse("failed", &orderID, createdOrder.Status, "", "La orden no esta en estado PENDING"))
		case errors.Is(err, service.ErrOrderNotFound):
			c.JSON(http.StatusNotFound, checkoutResponse("failed", &orderID, "", "", "No se encontro la orden recien creada"))
		default:
			c.JSON(http.StatusInternalServerError, checkoutResponse("failed", &orderID, createdOrder.Status, "", "No se pudo pagar la orden en checkout: "+err.Error()))
		}
		return
	}

	h.store.Del(key)

	paidOrder, err := h.orderService.GetOrderByID(orderID.String())
	if err != nil {
		c.JSON(http.StatusCreated, checkoutResponse("success", &orderID, model.OrderStatusPaid, "Checkout completado y orden pagada", ""))
		return
	}

	resp := checkoutResponse("success", &orderID, paidOrder.Status, "Checkout completado y orden pagada", "")
	resp["order"] = paidOrder
	c.JSON(http.StatusCreated, resp)
}

func (h *CartHandler) loadCart(key string, userID uuid.UUID) (model.Cart, bool, error) {
	data, exists := h.store.Get(key)
	if !exists {
		return model.Cart{UserID: userID}, false, nil
	}

	var cart model.Cart
	if err := json.Unmarshal(data, &cart); err != nil {
		return model.Cart{}, false, err
	}

	return cart, true, nil
}

func (h *CartHandler) saveCart(key string, cart model.Cart) error {
	cartBytes, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	h.store.Set(key, cartBytes, cartTTLSeconds)
	return nil
}

func cartKey(userID string) string {
	return fmt.Sprintf("cart:user:%s", userID)
}

func checkoutResponse(
	status string,
	orderID *uuid.UUID,
	orderStatus model.OrderStatus,
	message string,
	errorMessage string,
) gin.H {
	resp := gin.H{
		"status":       status,
		"order_id":     nil,
		"order_status": orderStatus,
	}

	if orderID != nil {
		resp["order_id"] = orderID.String()
	}

	if message != "" {
		resp["message"] = message
	}
	if errorMessage != "" {
		resp["error"] = errorMessage
	}

	return resp
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, string, bool) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return uuid.Nil, "", false
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario invalido en token"})
		return uuid.Nil, "", false
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID de usuario invalido"})
		return uuid.Nil, "", false
	}

	return userID, userIDStr, true
}

func isNotFoundError(err error) bool {
	errLower := strings.ToLower(err.Error())
	return strings.Contains(errLower, "no se encontro") || strings.Contains(errLower, "record not found")
}
