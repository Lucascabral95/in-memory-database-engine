package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucas-dev/in-memory-db/internal/config"
	"github.com/lucas-dev/in-memory-db/internal/storage"
)

func TestSetupRoutes_RegistersExpectedRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	cfg := &config.Config{
		Port:        "8080",
		Environment: "test",
		JWTSecret:   "test-secret",
	}

	SetupRoutes(router, cfg, nil, nil, nil, nil, storage.NewMemoryStore())

	expected := map[string]bool{
		"GET /health":                    false,
		"POST /users/register":           false,
		"POST /users/login":              false,
		"DELETE /users/:email":           false,
		"GET /categories":                false,
		"GET /categories/:id":            false,
		"POST /categories":               false,
		"PATCH /categories/:id":          false,
		"DELETE /categories/:id":         false,
		"GET /products":                  false,
		"GET /products/:id":              false,
		"POST /products":                 false,
		"PATCH /products/:id":            false,
		"DELETE /products/:id":           false,
		"GET /orders":                    false,
		"GET /orders/:id":                false,
		"POST /orders":                   false,
		"POST /orders/:id/pay":           false,
		"PATCH /orders/:id":              false,
		"POST /cart/items":               false,
		"PATCH /cart/items/:product_id":  false,
		"DELETE /cart/items/:product_id": false,
		"GET /cart":                      false,
		"DELETE /cart":                   false,
		"POST /cart/checkout":            false,
	}

	for _, route := range router.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := expected[key]; ok {
			expected[key] = true
		}
	}

	for route, found := range expected {
		if !found {
			t.Fatalf("expected route not registered: %s", route)
		}
	}
}

func TestSetupRoutes_HealthRequiresAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	cfg := &config.Config{
		Port:        "8080",
		Environment: "test",
		JWTSecret:   "test-secret",
	}

	SetupRoutes(router, cfg, nil, nil, nil, nil, storage.NewMemoryStore())

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("GET /health status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
