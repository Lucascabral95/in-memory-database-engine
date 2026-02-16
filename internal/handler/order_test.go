package handler

import (
	"net/http"
	"testing"
)

func TestOrderHandler_CreateOrder_Unauthorized(t *testing.T) {
	h := &OrderHandler{}
	ctx, rec := newUnauthorizedRequest(http.MethodPost, "/orders")

	h.CreateOrder(ctx)

	assertUnauthorized(t, rec)
}

func TestOrderHandler_GetAllOrders_Unauthorized(t *testing.T) {
	h := &OrderHandler{}
	ctx, rec := newUnauthorizedRequest(http.MethodGet, "/orders")

	h.GetAllOrders(ctx)

	assertUnauthorized(t, rec)
}

func TestOrderHandler_GetOrderByID_InvalidUUID(t *testing.T) {
	h := &OrderHandler{}
	ctx, rec := newJSONContext(http.MethodGet, "/orders/invalid", "")
	withParam(ctx, "id", "invalid-uuid")

	h.GetOrderByID(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestOrderHandler_UpdateStatusOrder_InvalidUUID(t *testing.T) {
	h := &OrderHandler{}
	ctx, rec := newJSONContext(http.MethodPatch, "/orders/invalid", "{}")
	withParam(ctx, "id", "invalid-uuid")

	h.UpdateStatusOrder(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
}

func TestOrderHandler_OrderUpdatePay_InvalidUUID(t *testing.T) {
	h := &OrderHandler{}
	ctx, rec := newJSONContext(http.MethodPost, "/orders/invalid/pay", "")
	withParam(ctx, "id", "invalid-uuid")

	h.OrderUpdatePay(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestOrderHandler_GetAuthUserID_Unauthorized(t *testing.T) {
	h := &OrderHandler{}
	ctx, rec := newJSONContext(http.MethodGet, "/", "")

	_, ok := h.getAuthUserID(ctx)
	if ok {
		t.Fatalf("getAuthUserID() ok = true, want false")
	}

	assertUnauthorized(t, rec)
}
