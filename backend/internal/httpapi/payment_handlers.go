package httpapi

import (
	"net/http"
)

// handleCreatePaymentIntent returns 201 with an empty PaymentIntent payload (zero values)
// Route: POST /orders/{orderId}/payment-intents (auth or X-Order-Secret via middleware)
func handleCreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	pi := PaymentIntentDTO{}
	writeJSON(w, http.StatusCreated, pi)
}

// handleGetPaymentIntent returns 200 with an empty PaymentIntent payload (zero values)
// Route: GET /orders/{orderId}/payment-intents/{paymentIntentId} (auth or X-Order-Secret)
func handleGetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	pi := PaymentIntentDTO{}
	writeJSON(w, http.StatusOK, pi)
}

// handlePaymentWebhook acknowledges provider webhook without auth.
// Route: POST /payments/webhook/paystack
func handlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	_ = r.Header.Get("X-Paystack-Signature") // read if present; do not verify
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}
