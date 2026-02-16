package utils

import "github.com/lucas-dev/in-memory-db/internal/model"

func IsValidOrderStatus(status model.OrderStatus) bool {
	switch status {
	case model.OrderStatusPending, model.OrderStatusPaid, model.OrderStatusShipped, model.OrderStatusCancelled:
		return true
	default:
		return false
	}
}
