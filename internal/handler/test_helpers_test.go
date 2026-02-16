package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newJSONContext(method, target, body string) (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rec)

	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	return ctx, rec
}

func withAuthUser(ctx *gin.Context) {
	ctx.Set("userID", uuid.NewString())
}

func withParam(ctx *gin.Context, key, value string) {
	ctx.Params = append(ctx.Params, gin.Param{Key: key, Value: value})
}

func decodeJSONBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	return body
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
}

func assertHasErrorKey(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	body := decodeJSONBody(t, rec)
	if _, ok := body["error"]; !ok {
		t.Fatalf("response body does not contain 'error' key: %#v", body)
	}
}

func newTestCartHandler() *CartHandler {
	return NewCartHandler(nil, nil, nil)
}

func newUnauthorizedRequest(method, target string) (*gin.Context, *httptest.ResponseRecorder) {
	ctx, rec := newJSONContext(method, target, "")
	ctx.Request = httptest.NewRequest(method, target, nil)
	return ctx, rec
}

func assertUnauthorized(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	assertStatus(t, rec.Code, http.StatusUnauthorized)
	assertHasErrorKey(t, rec)
}
