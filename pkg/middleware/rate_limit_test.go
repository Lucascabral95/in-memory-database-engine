package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func TestRateLimiterGetClient_ReusesLimiterForSameIP(t *testing.T) {
	rl := newRateLimiterWithConfig(rate.Every(time.Second), 2, time.Hour, time.Hour)

	l1 := rl.getClient("203.0.113.10")
	l2 := rl.getClient("203.0.113.10")

	if l1 != l2 {
		t.Fatalf("expected same limiter instance for same IP")
	}
	if len(rl.clients) != 1 {
		t.Fatalf("expected 1 tracked client, got %d", len(rl.clients))
	}
}

func TestRateLimiterMiddleware_AllowsThenBlocks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := newRateLimiterWithConfig(rate.Every(time.Hour), 1, time.Hour, time.Hour)
	r := gin.New()
	r.Use(rl.Middleware())
	r.GET("/limited", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req1 := httptest.NewRequest(http.MethodGet, "/limited", nil)
	req1.RemoteAddr = "198.51.100.7:12345"
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("first request status = %d, want %d", rec1.Code, http.StatusOK)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/limited", nil)
	req2.RemoteAddr = "198.51.100.7:12345"
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request status = %d, want %d", rec2.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimiterMiddleware_TooManyRequestsResponseBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := newRateLimiterWithConfig(rate.Every(time.Hour), 1, time.Hour, time.Hour)
	r := gin.New()
	r.Use(rl.Middleware())
	r.GET("/limited", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req1 := httptest.NewRequest(http.MethodGet, "/limited", nil)
	req1.RemoteAddr = "203.0.113.55:12345"
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)

	req2 := httptest.NewRequest(http.MethodGet, "/limited", nil)
	req2.RemoteAddr = "203.0.113.55:12345"
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec2.Code, http.StatusTooManyRequests)
	}

	var body map[string]string
	if err := json.Unmarshal(rec2.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json body: %v", err)
	}

	if body["error"] != "Too many requests" {
		t.Fatalf("error = %q, want %q", body["error"], "Too many requests")
	}
	if body["message"] != "Demasiadas solicitudes, intenta de nuevo mas tarde" {
		t.Fatalf("message = %q, want %q", body["message"], "Demasiadas solicitudes, intenta de nuevo mas tarde")
	}
}

func TestRateLimiterCleanup_RemovesInactiveClients(t *testing.T) {
	rl := newRateLimiterWithConfig(rate.Every(time.Second), 2, 10*time.Millisecond, 20*time.Millisecond)

	rl.mu.Lock()
	rl.clients["198.51.100.1"] = &client{
		limiter:  rate.NewLimiter(rate.Every(time.Second), 1),
		lastSeen: time.Now().Add(-time.Minute),
	}
	rl.mu.Unlock()

	time.Sleep(40 * time.Millisecond)

	rl.mu.Lock()
	defer rl.mu.Unlock()
	if _, ok := rl.clients["198.51.100.1"]; ok {
		t.Fatalf("expected inactive client to be removed by cleanup")
	}
}
