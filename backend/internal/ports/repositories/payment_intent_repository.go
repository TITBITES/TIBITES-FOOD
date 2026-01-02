package repositories

import "context"

type PaymentIntentRepository interface {
	Create(ctx context.Context, rec PaymentIntentRow) (PaymentIntentRow, error)
	GetByID(ctx context.Context, id string) (PaymentIntentRow, error)
	GetByOrder(ctx context.Context, orderID string) (PaymentIntentRow, error)
}

type PaymentIntentRow struct{
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
