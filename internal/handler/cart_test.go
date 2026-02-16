package handler

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestCartHandler_AddToCart_Unauthorized(t *testing.T) {
	h := newTestCartHandler()
	ctx, rec := newJSONContext(http.MethodPost, "/cart/items", `{"product_id":"x","quantity":1}`)

	h.AddToCart(ctx)

	assertUnauthorized(t, rec)
}

func TestCartHandler_AddToCart_InvalidJSON(t *testing.T) {
	h := newTestCartHandler()
	ctx, rec := newJSONContext(http.MethodPost, "/cart/items", "{")
	withAuthUser(ctx)

	h.AddToCart(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestCartHandler_AddToCart_InvalidQuantity(t *testing.T) {
	h := newTestCartHandler()
	ctx, rec := newJSONContext(http.MethodPost, "/cart/items", `{"product_id":"11111111-1111-1111-1111-111111111111","quantity":0}`)
	withAuthUser(ctx)

	h.AddToCart(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestCartHandler_UpdateCartItem_Unauthorized(t *testing.T) {
	h := newTestCartHandler()
	ctx, rec := newJSONContext(http.MethodPatch, "/cart/items/invalid", `{"quantity":1}`)
	withParam(ctx, "product_id", "invalid-uuid")

	h.UpdateCartItem(ctx)

	assertUnauthorized(t, rec)
}

func TestCartHandler_UpdateCartItem_InvalidProductID(t *testing.T) {
	h := newTestCartHandler()
	ctx, rec := newJSONContext(http.MethodPatch, "/cart/items/invalid", `{"quantity":1}`)
	withAuthUser(ctx)
	withParam(ctx, "product_id", "invalid-uuid")

	h.UpdateCartItem(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestCartHandler_RemoveCartItem_Unauthorized(t *testing.T) {
	h := newTestCartHandler()
	ctx, rec := newJSONContext(http.MethodDelete, "/cart/items/invalid", "")
	withParam(ctx, "product_id", "invalid-uuid")

	h.RemoveCartItem(ctx)

	assertUnauthorized(t, rec)
}

func TestCartHandler_RemoveCartItem_InvalidProductID(t *testing.T) {
	h := newTestCartHandler()
	ctx, rec := newJSONContext(http.MethodDelete, "/cart/items/invalid", "")
	withAuthUser(ctx)
	withParam(ctx, "product_id", "invalid-uuid")

	h.RemoveCartItem(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestCartHandler_ClearCart_Unauthorized(t *testing.T) {
	h := newTestCartHandler()
	ctx, rec := newJSONContext(http.MethodDelete, "/cart", "")

	h.ClearCart(ctx)

	assertUnauthorized(t, rec)
}

func TestCartHandler_GetCart_Unauthorized(t *testing.T) {
	h := newTestCartHandler()
	ctx, rec := newJSONContext(http.MethodGet, "/cart", "")

	h.GetCart(ctx)

	assertUnauthorized(t, rec)
}

func TestCartHandler_Checkout_Unauthorized(t *testing.T) {
	h := newTestCartHandler()
	ctx, rec := newJSONContext(http.MethodPost, "/cart/checkout", "")

	h.Checkout(ctx)

	assertUnauthorized(t, rec)
}

func TestCheckoutResponse_BuildsPayload(t *testing.T) {
	id := mustParseUUID(t, "11111111-1111-1111-1111-111111111111")
	resp := checkoutResponse("failed", &id, "", "", "boom")

	if resp["status"] != "failed" {
		t.Fatalf("status = %v, want %v", resp["status"], "failed")
	}
	if resp["order_id"] == nil {
		t.Fatalf("order_id should not be nil")
	}
	if resp["error"] != "boom" {
		t.Fatalf("error = %v, want %v", resp["error"], "boom")
	}
}

func mustParseUUID(t *testing.T, input string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(input)
	if err != nil {
		t.Fatalf("uuid.Parse() error: %v", err)
	}
	return id
}
