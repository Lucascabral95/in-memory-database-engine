package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lucas-dev/in-memory-db/pkg/utils"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusUnauthorized, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Token no proporcionado") {
		t.Errorf("Se esperaba error 'Token no proporcionado', se obtuvo: %s", w.Body.String())
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	tests := []struct {
		name          string
		authHeader    string
		expectedError string
	}{
		{
			name:          "Sin Bearer prefix",
			authHeader:    "invalid-token",
			expectedError: "Formato de token inválido",
		},
		{
			name:          "Solo Bearer",
			authHeader:    "Bearer",
			expectedError: "Formato de token inválido",
		},
		{
			name:          "Multiple partes",
			authHeader:    "Bearer token extra",
			expectedError: "Formato de token inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tt.authHeader)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusUnauthorized, w.Code)
			}
			if !strings.Contains(w.Body.String(), tt.expectedError) {
				t.Errorf("Se esperaba error '%s', se obtuvo: %s", tt.expectedError, w.Body.String())
			}
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router := gin.New()
	router.Use(AuthMiddleware("test-secret"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-jwt-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusUnauthorized, w.Code)
	}
	if !strings.Contains(w.Body.String(), "Token inválido o expirado") {
		t.Errorf("Se esperaba error 'Token inválido o expirado', se obtuvo: %s", w.Body.String())
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Setup
	jwtSecret := "test-secret"
	userID := "123e4567-e89b-12d3-a456-426614174000"

	// Generar token válido
	token, err := utils.GenerateToken(jwtSecret, userID, "test@example.com", "Test", "User")
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}
	if token == "" {
		t.Fatal("El token generado está vacío")
	}

	router := gin.New()
	router.Use(AuthMiddleware(jwtSecret))
	router.GET("/protected", func(c *gin.Context) {
		// Verificamos que los valores se setearon en el contexto
		userIDValue, exists := c.Get("userID")
		if !exists {
			t.Error("userID no encontrado en el contexto")
		}
		if userID != userIDValue {
			t.Errorf("Se esperaba userID %s, se obtuvo %v", userID, userIDValue)
		}

		email, exists := c.Get("userEmail")
		if !exists {
			t.Error("userEmail no encontrado en el contexto")
		}
		if "test@example.com" != email {
			t.Errorf("Se esperaba email %s, se obtuvo %v", "test@example.com", email)
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), "ok") {
		t.Errorf("Se esperaba cuerpo 'ok', se obtuvo: %s", w.Body.String())
	}
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	token, err := utils.GenerateToken("secret-A", "user-id", "test@example.com", "Test", "User")
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}

	router := gin.New()
	router.Use(AuthMiddleware("secret-B"))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusUnauthorized, w.Code)
	}
	if !strings.Contains(w.Body.String(), "Token inválido o expirado") {
		t.Errorf("Se esperaba error 'Token inválido o expirado', se obtuvo: %s", w.Body.String())
	}
}

func TestAuthMiddleware_ContextValues(t *testing.T) {
	jwtSecret := "test-secret"
	userID := "123e4567-e89b-12d3-a456-426614174000"
	email := "test@example.com"
	firstName := "Juan"
	lastName := "Perez"

	token, err := utils.GenerateToken(jwtSecret, userID, email, firstName, lastName)
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}

	router := gin.New()
	router.Use(AuthMiddleware(jwtSecret))
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"userID":        c.GetString("userID"),
			"userEmail":     c.GetString("userEmail"),
			"userFirstName": c.GetString("userFirstName"),
			"userLastName":  c.GetString("userLastName"),
		})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, userID) {
		t.Errorf("El cuerpo no contiene el userID: %s", body)
	}
	if !strings.Contains(body, email) {
		t.Errorf("El cuerpo no contiene el email: %s", body)
	}
	if !strings.Contains(body, firstName) {
		t.Errorf("El cuerpo no contiene el firstName: %s", body)
	}
	if !strings.Contains(body, lastName) {
		t.Errorf("El cuerpo no contiene el lastName: %s", body)
	}
}
