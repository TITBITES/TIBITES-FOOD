package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	porttokens "local.dev/foodapp/internal/ports/tokens"
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
func NewRouter(tokenSvc porttokens.TokenService) http.Handler {
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
		rm.Use(requireAuthWith(tokenSvc))
		rm.Get("/users/me", func(w http.ResponseWriter, r *http.Request) {
			notImplemented(w, r, "NOT_IMPLEMENTED", "users/me not implemented")
		})
	})

	// Menu (read-only)
	r.Get("/menu/categories", handleListCategories)
	r.Get("/menu/items", handleListMenuItems)

	// Orders
	r.Group(func(rm chi.Router) {
		// GET /orders requires auth
		rm.With(requireAuthWith(tokenSvc)).Get("/orders", handleListOrders)

		// POST /orders (guest allowed)
		rm.Post("/orders", handleCreateOrder)

		// GET /orders/{orderId} (auth or order secret)
		rm.With(authOrOrderSecretWith(tokenSvc)).Get("/orders/{orderId}", handleGetOrder)
	})

	// Payments
	r.Group(func(rm chi.Router) {
		// Create PI (auth or order secret)
		rm.With(authOrOrderSecretWith(tokenSvc)).Post("/orders/{orderId}/payment-intents", handleCreatePaymentIntent)
		// Get PI (auth or order secret)
		rm.With(authOrOrderSecretWith(tokenSvc)).Get("/orders/{orderId}/payment-intents/{paymentIntentId}", handleGetPaymentIntent)
		// Webhook (no auth)
		rm.Post("/payments/webhook/paystack", handlePaymentWebhook)
	})

	return r
}

// requireAuth validates Authorization: Bearer <token> using TokenService.
func requireAuthWith(tokenSvc porttokens.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				writeJSON(w, http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "missing Authorization header"})
				return
			}
			if !strings.HasPrefix(auth, "Bearer ") {
				writeJSON(w, http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "invalid Authorization header"})
				return
			}
			tok := strings.TrimPrefix(auth, "Bearer ")
			if _, _, err := tokenSvc.ValidateAccessToken(r.Context(), tok); err != nil {
				writeJSON(w, http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "invalid or expired token"})
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
func authOrOrderSecretWith(tokenSvc porttokens.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			secret := r.Header.Get("X-Order-Secret")
			if auth == "" && secret == "" {
				writeJSON(w, http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "missing Authorization or X-Order-Secret"})
				return
			}
			if auth != "" && strings.HasPrefix(auth, "Bearer ") {
				tok := strings.TrimPrefix(auth, "Bearer ")
				if _, _, err := tokenSvc.ValidateAccessToken(r.Context(), tok); err != nil {
					writeJSON(w, http.StatusUnauthorized, Error{Code: "UNAUTHORIZED", Message: "invalid or expired token"})
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
