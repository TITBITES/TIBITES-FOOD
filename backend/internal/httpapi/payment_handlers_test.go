package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"local.dev/foodapp/internal/adapters/tokens"
)

func TestPaymentIntentRoutesAuth(t *testing.T) {
	tokenSvc := tokens.NewJWTService("secret1", "secret2", time.Minute, time.Hour*24)
	r := NewRouter(tokenSvc)

	// POST create PI requires auth or order secret
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/orders/00000000-0000-0000-0000-000000000000/payment-intents", nil)
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for POST payment-intents without auth or secret, got %d", rw.Code)
	}

	// with X-Order-Secret should be allowed
	rw2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/orders/00000000-0000-0000-0000-000000000000/payment-intents", nil)
	req2.Header.Set("X-Order-Secret", "dummy")
	r.ServeHTTP(rw2, req2)
	if rw2.Code != http.StatusCreated {
		t.Fatalf("expected 201 for POST payment-intents with order secret, got %d", rw2.Code)
	}

	// GET PI requires auth or order secret
	rw3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/orders/00000000-0000-0000-0000-000000000000/payment-intents/11111111-1111-1111-1111-111111111111", nil)
	r.ServeHTTP(rw3, req3)
	if rw3.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for GET payment-intents without auth or secret, got %d", rw3.Code)
	}

	// with X-Order-Secret should be allowed
	rw4 := httptest.NewRecorder()
	req4 := httptest.NewRequest(http.MethodGet, "/orders/00000000-0000-0000-0000-000000000000/payment-intents/11111111-1111-1111-1111-111111111111", nil)
	req4.Header.Set("X-Order-Secret", "dummy")
	r.ServeHTTP(rw4, req4)
	if rw4.Code != http.StatusOK {
		t.Fatalf("expected 200 for GET payment-intents with order secret, got %d", rw4.Code)
	}
}

func TestPaymentsWebhookBypassesAuth(t *testing.T) {
	tokenSvc := tokens.NewJWTService("secret1", "secret2", time.Minute, time.Hour*24)
	r := NewRouter(tokenSvc)
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/payments/webhook/paystack", nil)
	// No auth headers
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 for webhook without auth, got %d", rw.Code)
	}
}
