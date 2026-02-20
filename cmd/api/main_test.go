package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestServerAddr(t *testing.T) {
	got := serverAddr("8080")
	want := ":8080"
	if got != want {
		t.Fatalf("serverAddr() = %q, want %q", got, want)
	}
}

func TestRegisterSwaggerRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	registerSwaggerRoute(r)

	req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Fatalf("GET /swagger/index.html = 404, want non-404")
	}
}
