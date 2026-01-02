package payments

import "context"

type ProviderAdapter interface {
	CreatePaymentIntent(ctx context.Context, order ProviderOrderSnapshot, options map[string]any) (ProviderPaymentIntent, error)
	GetPaymentIntent(ctx context.Context, providerIntentID string) (ProviderPaymentIntent, error)
}

type WebhookVerifier interface {
	VerifySignature(headers map[string]string, body []byte) (bool, error)
}

type ProviderOrderSnapshot struct{
	OrderID string
	AmountCents int
	Currency string
}

type ProviderPaymentIntent struct{
	ID string
	Status string
	ClientSecret *string
	AuthorizationURL *string
}
