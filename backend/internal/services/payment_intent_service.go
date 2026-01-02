package services

import "context"

type PaymentIntentService interface {
	CreateForOrder(ctx context.Context, orderID string, auth AuthContext, req PaymentIntentCreateRequest) (PaymentIntentDTO, error)
	Get(ctx context.Context, orderID, paymentIntentID string, auth AuthContext) (PaymentIntentDTO, error)
}

type PaymentIntentCreateRequest struct{
	ReturnURL *string
	ProviderOptions map[string]any
}

type PaymentIntentDTO struct{
	ID string
	OrderID string
	Provider string
	ProviderIntentID string
	Status string
	ClientSecret *string
	AuthorizationURL *string
	AmountCents int
	Currency string
	CreatedAt string
	UpdatedAt string
}
