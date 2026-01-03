package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"local.dev/foodapp/internal/adapters/tokens"
)

func TestOrdersRoutesExistAndStatuses(t *testing.T) {
	tokenSvc := tokens.NewJWTService("secret1", "secret2", time.Minute, time.Hour*24)
	r := NewRouter(tokenSvc)

	// POST /orders (guest allowed)
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/orders", nil)
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatalf("expected 201 for POST /orders, got %d", rw.Code)
	}

	// GET /orders requires auth
	rw2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/orders", nil)
	r.ServeHTTP(rw2, req2)
	if rw2.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for GET /orders without auth, got %d", rw2.Code)
	}

	// GET /orders/{id} allows order secret
	rw3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/orders/00000000-0000-0000-0000-000000000000", nil)
	r.ServeHTTP(rw3, req3)
	if rw3.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for GET /orders/{id} without auth or secret, got %d", rw3.Code)
	}

	// Now with X-Order-Secret
	rw4 := httptest.NewRecorder()
	req4 := httptest.NewRequest(http.MethodGet, "/orders/00000000-0000-0000-0000-000000000000", nil)
	req4.Header.Set("X-Order-Secret", "dummy")
	r.ServeHTTP(rw4, req4)
	if rw4.Code != http.StatusOK {
		t.Fatalf("expected 200 for GET /orders/{id} with X-Order-Secret, got %d", rw4.Code)
	}
}
