package handler

import (
	"net/http"
	"testing"
)

func TestCategoryHandler_CreateCategory_InvalidJSON(t *testing.T) {
	h := &CategoryHandler{}
	ctx, rec := newJSONContext(http.MethodPost, "/categories", "{")

	h.CreateCategory(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestCategoryHandler_FindCategoryByID_InvalidUUID(t *testing.T) {
	h := &CategoryHandler{}
	ctx, rec := newJSONContext(http.MethodGet, "/categories/invalid", "")
	withParam(ctx, "id", "invalid-uuid")

	h.FindCategoryByID(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestCategoryHandler_UpdateCategory_InvalidUUID(t *testing.T) {
	h := &CategoryHandler{}
	ctx, rec := newJSONContext(http.MethodPatch, "/categories/invalid", "{}")
	withParam(ctx, "id", "invalid-uuid")

	h.UpdateCategory(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestCategoryHandler_DeleteCategory_InvalidUUID(t *testing.T) {
	h := &CategoryHandler{}
	ctx, rec := newJSONContext(http.MethodDelete, "/categories/invalid", "")
	withParam(ctx, "id", "invalid-uuid")

	h.DeleteCategory(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}
