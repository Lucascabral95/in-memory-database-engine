package handler

import (
	"net/http"
	"testing"
)

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	h := &ProductHandler{}
	ctx, rec := newJSONContext(http.MethodPost, "/products", "{")

	h.CreateProduct(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestProductHandler_GetProductByID_InvalidUUID(t *testing.T) {
	h := &ProductHandler{}
	ctx, rec := newJSONContext(http.MethodGet, "/products/invalid", "")
	withParam(ctx, "id", "invalid-uuid")

	h.GetProductByID(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestProductHandler_UpdateProduct_InvalidUUID(t *testing.T) {
	h := &ProductHandler{}
	ctx, rec := newJSONContext(http.MethodPatch, "/products/invalid", "{}")
	withParam(ctx, "id", "invalid-uuid")

	h.UpdateProduct(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestProductHandler_DeleteProduct_InvalidUUID(t *testing.T) {
	h := &ProductHandler{}
	ctx, rec := newJSONContext(http.MethodDelete, "/products/invalid", "")
	withParam(ctx, "id", "invalid-uuid")

	h.DeleteProduct(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}
