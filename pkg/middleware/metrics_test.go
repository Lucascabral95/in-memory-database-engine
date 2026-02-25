package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMetricsHandler_ExposesPrometheusMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(MetricsMiddleware())
	r.GET("/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.GET("/metrics", MetricsHandler())

	reqProducts := httptest.NewRequest(http.MethodGet, "/products", nil)
	recProducts := httptest.NewRecorder()
	r.ServeHTTP(recProducts, reqProducts)

	if recProducts.Code != http.StatusOK {
		t.Fatalf("GET /products status = %d, want %d", recProducts.Code, http.StatusOK)
	}

	reqMetrics := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	recMetrics := httptest.NewRecorder()
	r.ServeHTTP(recMetrics, reqMetrics)

	if recMetrics.Code != http.StatusOK {
		t.Fatalf("GET /metrics status = %d, want %d", recMetrics.Code, http.StatusOK)
	}

	body := recMetrics.Body.String()
	if !strings.Contains(body, "http_requests_total") {
		t.Fatalf("expected /metrics to include http_requests_total")
	}
	if !strings.Contains(body, "http_request_duration_seconds") {
		t.Fatalf("expected /metrics to include http_request_duration_seconds")
	}
	if !strings.Contains(body, `path="/products"`) {
		t.Fatalf("expected /metrics to include /products labels")
	}
}
