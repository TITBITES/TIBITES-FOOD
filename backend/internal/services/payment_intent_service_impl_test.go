package services

import (
	"context"
	"testing"

	paystub "local.dev/foodapp/internal/adapters/payments/stub"
	repo "local.dev/foodapp/internal/adapters/repositories/memory"
)

func newPIService() (PaymentIntentService, OrderService) {
	orders := repo.NewOrderRepoMemory()
	secrets := repo.NewOrderSecretRepoMemory()
	pis := repo.NewPaymentIntentRepoMemory()
	provider := paystub.NewProviderStub()
	piSvc := NewPaymentIntentService(orders, secrets, pis, provider)
	ordSvc := NewOrderService(orders, secrets)
	return piSvc, ordSvc
}

func TestCreatePIAsGuestWithSecret(t *testing.T) {
	piSvc, ordSvc := newPIService()
	ctx := context.Background()
	order, secret, _ := ordSvc.CreateOrderGuest(ctx, OrderCreateRequest{Currency: "NGN", Items: []OrderItemCreate{{MenuItemID: "item-1", Quantity: 1}}})
	_ = secret
	pi, err := piSvc.CreateForOrder(ctx, order.ID, AuthContext{OrderSecret: &secret}, PaymentIntentCreateRequest{})
	if err != nil { t.Fatalf("err: %v", err) }
	if pi.ID == "" || pi.ProviderIntentID == "" { t.Fatalf("expected IDs") }
}

func TestCreatePIAsCustomer(t *testing.T) {
	piSvc, ordSvc := newPIService()
	ctx := context.Background()
	order, _ := ordSvc.CreateOrderCustomer(ctx, "cust-1", OrderCreateRequest{Currency: "NGN", Items: []OrderItemCreate{{MenuItemID: "item-2", Quantity: 1}}})
	pi, err := piSvc.CreateForOrder(ctx, order.ID, AuthContext{CustomerID: strPtr("cust-1")}, PaymentIntentCreateRequest{})
	if err != nil { t.Fatalf("err: %v", err) }
	if pi.ID == "" { t.Fatalf("expected id") }
}

func TestCreatePIReturnsExisting(t *testing.T) {
	piSvc, ordSvc := newPIService()
	ctx := context.Background()
	order, secret, _ := ordSvc.CreateOrderGuest(ctx, OrderCreateRequest{Currency: "NGN", Items: []OrderItemCreate{{MenuItemID: "item-1", Quantity: 1}}})
	first, _ := piSvc.CreateForOrder(ctx, order.ID, AuthContext{OrderSecret: &secret}, PaymentIntentCreateRequest{})
	second, _ := piSvc.CreateForOrder(ctx, order.ID, AuthContext{OrderSecret: &secret}, PaymentIntentCreateRequest{})
	if first.ID != second.ID { t.Fatalf("expected same PI returned") }
}

func TestGetPIWithSecretOrCustomer(t *testing.T) {
	piSvc, ordSvc := newPIService()
	ctx := context.Background()
	order, secret, _ := ordSvc.CreateOrderGuest(ctx, OrderCreateRequest{Currency: "NGN", Items: []OrderItemCreate{{MenuItemID: "item-1", Quantity: 1}}})
	created, _ := piSvc.CreateForOrder(ctx, order.ID, AuthContext{OrderSecret: &secret}, PaymentIntentCreateRequest{})
	// via secret
	_, err := piSvc.Get(ctx, order.ID, created.ID, AuthContext{OrderSecret: &secret})
	if err != nil { t.Fatalf("secret get failed: %v", err) }
	// via customer
	order2, _ := ordSvc.CreateOrderCustomer(ctx, "cust-1", OrderCreateRequest{Currency: "NGN", Items: []OrderItemCreate{{MenuItemID: "item-2", Quantity: 1}}})
	created2, _ := piSvc.CreateForOrder(ctx, order2.ID, AuthContext{CustomerID: strPtr("cust-1")}, PaymentIntentCreateRequest{})
	_, err = piSvc.Get(ctx, order2.ID, created2.ID, AuthContext{CustomerID: strPtr("cust-1")})
	if err != nil { t.Fatalf("customer get failed: %v", err) }
}

func TestCannotCreatePIForForeignOrder(t *testing.T) {
	piSvc, ordSvc := newPIService()
	ctx := context.Background()
	order, secret, _ := ordSvc.CreateOrderGuest(ctx, OrderCreateRequest{Currency: "NGN", Items: []OrderItemCreate{{MenuItemID: "item-1", Quantity: 1}}})
	_ = secret
	_, err := piSvc.CreateForOrder(ctx, order.ID, AuthContext{CustomerID: strPtr("cust-x")}, PaymentIntentCreateRequest{})
	if err == nil { t.Fatalf("expected forbidden error") }
	_, err = piSvc.CreateForOrder(ctx, order.ID, AuthContext{OrderSecret: strPtr("wrong")}, PaymentIntentCreateRequest{})
	if err == nil { t.Fatalf("expected forbidden error") }
}

func strPtr(s string) *string { return &s }
