package services

import (
	"context"
	"errors"

	"local.dev/foodapp/internal/ports/payments"
	"local.dev/foodapp/internal/ports/repositories"
)

type paymentIntentService struct{
	orders repositories.OrderRepository
	secrets repositories.OrderSecretRepository
	pis repositories.PaymentIntentRepository
	provider payments.ProviderAdapter
}

func NewPaymentIntentService(orders repositories.OrderRepository, secrets repositories.OrderSecretRepository, pis repositories.PaymentIntentRepository, provider payments.ProviderAdapter) PaymentIntentService {
	return &paymentIntentService{orders: orders, secrets: secrets, pis: pis, provider: provider}
}

func (s *paymentIntentService) CreateForOrder(ctx context.Context, orderID string, auth AuthContext, req PaymentIntentCreateRequest) (PaymentIntentDTO, error) {
	// Validate order exists
	order, _, err := s.orders.GetByID(ctx, orderID)
	if err != nil { return PaymentIntentDTO{}, err }
	// Access control
	if auth.CustomerID != nil {
		if order.CustomerID == nil || *order.CustomerID != *auth.CustomerID { return PaymentIntentDTO{}, errors.New("forbidden") }
	} else if auth.OrderSecret != nil {
		ok, err := s.secrets.Validate(ctx, orderID, *auth.OrderSecret)
		if err != nil { return PaymentIntentDTO{}, err }
		if !ok { return PaymentIntentDTO{}, errors.New("forbidden") }
	} else {
		return PaymentIntentDTO{}, errors.New("unauthorized")
	}
	// One active per order: return existing if present
	if existing, err := s.pis.GetByOrder(ctx, orderID); err == nil {
		return mapPI(existing), nil
	}
	// Call provider stub
	provPI, err := s.provider.CreatePaymentIntent(ctx, payments.ProviderOrderSnapshot{OrderID: orderID, AmountCents: order.TotalCents, Currency: order.Currency}, req.ProviderOptions)
	if err != nil { return PaymentIntentDTO{}, err }
	// Persist
	rec, err := s.pis.Create(ctx, repositories.PaymentIntentRow{
		OrderID: orderID,
		Provider: "stub",
		ProviderIntentID: provPI.ID,
		Status: provPI.Status,
		ClientSecret: provPI.ClientSecret,
		AuthorizationURL: provPI.AuthorizationURL,
		AmountCents: order.TotalCents,
		Currency: order.Currency,
	})
	if err != nil { return PaymentIntentDTO{}, err }
	return mapPI(rec), nil
}

func (s *paymentIntentService) Get(ctx context.Context, orderID, paymentIntentID string, auth AuthContext) (PaymentIntentDTO, error) {
	// Validate order exists
	order, _, err := s.orders.GetByID(ctx, orderID)
	if err != nil { return PaymentIntentDTO{}, err }
	// Access
	if auth.CustomerID != nil {
		if order.CustomerID == nil || *order.CustomerID != *auth.CustomerID { return PaymentIntentDTO{}, errors.New("forbidden") }
	} else if auth.OrderSecret != nil {
		ok, err := s.secrets.Validate(ctx, orderID, *auth.OrderSecret)
		if err != nil { return PaymentIntentDTO{}, err }
		if !ok { return PaymentIntentDTO{}, errors.New("forbidden") }
	} else { return PaymentIntentDTO{}, errors.New("unauthorized") }
	// Fetch PI
	rec, err := s.pis.GetByID(ctx, paymentIntentID)
	if err != nil { return PaymentIntentDTO{}, err }
	if rec.OrderID != orderID { return PaymentIntentDTO{}, errors.New("not found") }
	return mapPI(rec), nil
}

func mapPI(r repositories.PaymentIntentRow) PaymentIntentDTO {
	return PaymentIntentDTO{
		ID: r.ID,
		OrderID: r.OrderID,
		Provider: r.Provider,
		ProviderIntentID: r.ProviderIntentID,
		Status: r.Status,
		ClientSecret: r.ClientSecret,
		AuthorizationURL: r.AuthorizationURL,
		AmountCents: r.AmountCents,
		Currency: r.Currency,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
