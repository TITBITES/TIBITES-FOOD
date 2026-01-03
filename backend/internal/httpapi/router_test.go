package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"local.dev/foodapp/internal/adapters/tokens"
)

func TestAuthRequiredEndpoints(t *testing.T) {
	tokenSvc := tokens.NewJWTService("secret1", "secret2", time.Minute, time.Hour*24)
	r := NewRouter(tokenSvc)

	// GET /orders requires auth
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for GET /orders without auth, got %d", rw.Code)
	}

	// GET /orders/{orderId} requires auth or order secret
	req2 := httptest.NewRequest(http.MethodGet, "/orders/00000000-0000-0000-0000-000000000000", nil)
	rw2 := httptest.NewRecorder()
	r.ServeHTTP(rw2, req2)
	if rw2.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for GET /orders/{id} without auth or secret, got %d", rw2.Code)
	}

	// POST /orders/{orderId}/payment-intents requires auth or order secret
	req3 := httptest.NewRequest(http.MethodPost, "/orders/00000000-0000-0000-0000-000000000000/payment-intents", nil)
	rw3 := httptest.NewRecorder()
	r.ServeHTTP(rw3, req3)
	if rw3.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for POST payment-intents without auth or secret, got %d", rw3.Code)
	}
}

func TestMenuEndpoints(t *testing.T) {
	tokenSvc := tokens.NewJWTService("secret1", "secret2", time.Minute, time.Hour*24)
	r := NewRouter(tokenSvc)
	req := httptest.NewRequest(http.MethodGet, "/menu/categories", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 for GET /menu/categories, got %d", rw.Code)
	}
}
