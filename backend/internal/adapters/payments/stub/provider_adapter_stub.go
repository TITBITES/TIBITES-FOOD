package stub

import (
	"context"
	"strconv"

	"local.dev/foodapp/internal/ports/payments"
)

type ProviderStub struct{ seq int }

func NewProviderStub() *ProviderStub { return &ProviderStub{} }

func (p *ProviderStub) nextID() string { p.seq++; return "prov-" + strconv.Itoa(p.seq) }

func (p *ProviderStub) CreatePaymentIntent(ctx context.Context, order payments.ProviderOrderSnapshot, options map[string]any) (payments.ProviderPaymentIntent, error) {
	id := p.nextID()
	pi := payments.ProviderPaymentIntent{
		ID: id,
		Status: "pending",
		ClientSecret: strPtr("stub_client_secret"),
		AuthorizationURL: strPtr("https://stub.pay/authorize"),
	}
	return pi, nil
}

func (p *ProviderStub) GetPaymentIntent(ctx context.Context, providerIntentID string) (payments.ProviderPaymentIntent, error) {
	return payments.ProviderPaymentIntent{ID: providerIntentID, Status: "pending"}, nil
}

func strPtr(s string) *string { return &s }
