package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type contextKey string

const (
	ctxKeyUserID      contextKey = "userID"
	ctxKeyOrderSecret contextKey = "orderSecret"
)

// Error response matches docs/components.schemas.Error
type Error struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func notImplemented(w http.ResponseWriter, _ *http.Request, code, message string) {
	writeJSON(w, http.StatusNotImplemented, Error{Code: code, Message: message})
}

// NewRouter wires routes per OpenAPI 0.2.0. No business logic.
func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Auth
	r.Post("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		notImplemented(w, r, "NOT_IMPLEMENTED", "register not implemented")
	})
	r.Post("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		notImplemented(w, r, "NOT_IMPLEMENTED", "login not implemented")
	})
	r.Post("/auth/refresh", func(w http.ResponseWriter, r *http.Request) {
		notImplemented(w, r, "NOT_IMPLEMENTED", "refresh not implemented")
	})
	// Admin OAuth2
	r.Get("/auth/admin/oauth2/start", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"authorizationUrl": "https://example.com/oauth2/authorize"})
	})
	r.Get("/auth/admin/oauth2/callback", func(w http.ResponseWriter, r *http.Request) {
		notImplemented(w, r, "NOT_IMPLEMENTED", "oauth2 callback not implemented")
	})

	// Users
	r.Group(func(rm chi.Router) {
		rm.Use(requireAuth())
		rm.Get("/users/me", func(w http.ResponseWriter, r *http.Request) {
			notImplemented(w, r, "NOT_IMPLEMENTED", "users/me not implemented")
		})
	})

	// Menu (read-only)
	r.Get("/menu/categories", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{"categories": []interface{}{}})
	})
	r.Get("/menu/items", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"items":      []interface{}{},
			"pagination": map[string]int{"page": 1, "pageSize": 0, "total": 0},
		})
	})

	// Orders
	r.Group(func(rm chi.Router) {
		// GET /orders requires auth
		rm.With(requireAuth()).Get("/orders", func(w http.ResponseWriter, r *http.Request) {
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"orders":     []interface{}{},
				"pagination": map[string]int{"page": 1, "pageSize": 0, "total": 0},
			})
		})

		// POST /orders (guest allowed)
		rm.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
			notImplemented(w, r, "NOT_IMPLEMENTED", "create order not implemented")
		})

		// GET /orders/{orderId} (auth or order secret)
		rm.With(authOrOrderSecret()).Get("/orders/{orderId}", func(w http.ResponseWriter, r *http.Request) {
			notImplemented(w, r, "NOT_IMPLEMENTED", "get order not implemented")
		})
	})

	// Payments
	r.Group(func(rm chi.Router) {
		// Create PI (auth or order secret)
		rm.With(authOrOrderSecret()).Post("/orders/{orderId}/payment-intents", func(w http.ResponseWriter, r *http.Request) {
			notImplemented(w, r, "NOT_IMPLEMENTED", "create payment intent not implemented")
		})
		// Get PI (auth or order secret)
		rm.With(authOrOrderSecret()).Get("/orders/{orderId}/payment-intents/{paymentIntentId}", func(w http.ResponseWriter, r *http.Request) {
			notImplemented(w, r, "NOT_IMPLEMENTED", "get payment intent not implemented")
		})
		// Webhook (no auth)
		rm.Post("/payments/webhook/paystack", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		})
	})

	return r
}

// requireAuth enforces presence of Authorization: Bearer <token> header only (no verification).
func requireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				writeJSON(w, http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "missing Authorization header"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// orderSecretMiddleware extracts X-Order-Secret if present.
func orderSecretMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := r.Header.Get("X-Order-Secret")
			if secret != "" {
				// Attach to context if needed in future
				r = r.WithContext(r.Context())
			}
			next.ServeHTTP(w, r)
		})
	}
}

// authOrOrderSecret allows either Authorization or X-Order-Secret. If neither, returns 401.
func authOrOrderSecret() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			secret := r.Header.Get("X-Order-Secret")
			if auth == "" && secret == "" {
				writeJSON(w, http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "missing Authorization or X-Order-Secret"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
