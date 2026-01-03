package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"local.dev/foodapp/internal/adapters/tokens"
)

func TestHandleListCategories(t *testing.T) {
	tokenSvc := tokens.NewJWTService("secret1", "secret2", time.Minute, time.Hour*24)
	r := NewRouter(tokenSvc)
	req := httptest.NewRequest(http.MethodGet, "/menu/categories", nil)
	rw := httptest.NewRecorder()

	// Should not panic and should return 200
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestHandleListMenuItems_Defaults(t *testing.T) {
	tokenSvc := tokens.NewJWTService("secret1", "secret2", time.Minute, time.Hour*24)
	r := NewRouter(tokenSvc)
	req := httptest.NewRequest(http.MethodGet, "/menu/items", nil)
	rw := httptest.NewRecorder()

	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}

func TestHandleListMenuItems_WithParams(t *testing.T) {
	tokenSvc := tokens.NewJWTService("secret1", "secret2", time.Minute, time.Hour*24)
	r := NewRouter(tokenSvc)
	req := httptest.NewRequest(http.MethodGet, "/menu/items?page=2&pageSize=25&categoryId=00000000-0000-0000-0000-000000000000", nil)
	rw := httptest.NewRecorder()

	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
}
