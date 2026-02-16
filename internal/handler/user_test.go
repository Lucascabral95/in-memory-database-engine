package handler

import (
	"net/http"
	"testing"
)

func TestUserHandler_RegisterUser_InvalidJSON(t *testing.T) {
	h := &UserHandler{}
	ctx, rec := newJSONContext(http.MethodPost, "/users/register", "{")

	h.RegisterUser(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestUserHandler_LoginUser_InvalidJSON(t *testing.T) {
	h := &UserHandler{}
	ctx, rec := newJSONContext(http.MethodPost, "/users/login", "{")

	h.LoginUser(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestUserHandler_UpdateUser_InvalidJSON(t *testing.T) {
	h := &UserHandler{}
	ctx, rec := newJSONContext(http.MethodPatch, "/users", "{")

	h.UpdateUser(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}

func TestUserHandler_DeleteUser_EmptyEmail(t *testing.T) {
	h := &UserHandler{}
	ctx, rec := newJSONContext(http.MethodDelete, "/users/", "")
	withParam(ctx, "email", "")

	h.DeleteUser(ctx)

	assertStatus(t, rec.Code, http.StatusBadRequest)
	assertHasErrorKey(t, rec)
}
