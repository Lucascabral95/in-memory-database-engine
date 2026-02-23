package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetCorsConfig_SetsCORSHeadersOnSimpleRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(GetCorsConfig())
	r.GET("/resource", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/resource", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin == "" {
		t.Fatalf("expected Access-Control-Allow-Origin header to be set")
	}
}

func TestGetCorsConfig_PreflightIncludesConfiguredMethodsAndHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(GetCorsConfig())
	r.POST("/resource", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodOptions, "/resource", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "Authorization, Content-Type, X-Internal-Secret")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	allowMethods := rec.Header().Get("Access-Control-Allow-Methods")
	if !strings.Contains(allowMethods, http.MethodPost) || !strings.Contains(allowMethods, http.MethodPatch) {
		t.Fatalf("unexpected Access-Control-Allow-Methods: %q", allowMethods)
	}

	allowHeaders := rec.Header().Get("Access-Control-Allow-Headers")
	if !strings.Contains(strings.ToLower(allowHeaders), "authorization") {
		t.Fatalf("expected Authorization in Access-Control-Allow-Headers, got %q", allowHeaders)
	}
	if !strings.Contains(strings.ToLower(allowHeaders), "x-internal-secret") {
		t.Fatalf("expected X-Internal-Secret in Access-Control-Allow-Headers, got %q", allowHeaders)
	}

	if rec.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatalf("expected Access-Control-Allow-Credentials=true")
	}
}
