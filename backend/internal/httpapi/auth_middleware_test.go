package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"local.dev/foodapp/internal/adapters/tokens"
)

func TestMiddlewareRejectsInvalidToken(t *testing.T) {
	tokenSvc := tokens.NewJWTService("secret1", "secret2", time.Minute, time.Hour)
	r := NewRouter(tokenSvc)
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized { t.Fatalf("expected 401, got %d", rw.Code) }
}

func TestMiddlewareAcceptsValidToken(t *testing.T) {
	tokenSvc := tokens.NewJWTService("secret1", "secret2", time.Minute, time.Hour)
	r := NewRouter(tokenSvc)
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	access, _ := tokenSvc.GenerateAccessToken(req.Context(), "user-1", "Customer")
	req.Header.Set("Authorization", "Bearer "+access)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code == http.StatusUnauthorized { t.Fatalf("did not expect 401 with valid token") }
}
