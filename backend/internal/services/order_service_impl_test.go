package services

import (
	"context"
	"testing"

	repo "local.dev/foodapp/internal/adapters/repositories/memory"
)

func newOrderSvc() OrderService {
	orders := repo.NewOrderRepoMemory()
	secrets := repo.NewOrderSecretRepoMemory()
	return NewOrderService(orders, secrets)
}

func TestCreateGuestOrderReturnsSecret(t *testing.T) {
	svc := newOrderSvc()
	ctx := context.Background()
	dto, secret, err := svc.CreateOrderGuest(ctx, OrderCreateRequest{Currency: "NGN", TipCents: 0, Items: []OrderItemCreate{{MenuItemID: "item-1", Quantity: 1}}})
	if err != nil { t.Fatalf("err: %v", err) }
	if dto.ID == "" || secret == "" { t.Fatalf("expected id and secret") }
}

func TestCreateAuthOrderNoSecret(t *testing.T) {
	svc := newOrderSvc()
	ctx := context.Background()
	dto, err := svc.CreateOrderCustomer(ctx, "cust-1", OrderCreateRequest{Currency: "NGN", TipCents: 0, Items: []OrderItemCreate{{MenuItemID: "item-2", Quantity: 2}}})
	if err != nil { t.Fatalf("err: %v", err) }
	if dto.ID == "" { t.Fatalf("expected id") }
}

func TestGetOrderByCustomerOrSecret(t *testing.T) {
	svc := newOrderSvc()
	ctx := context.Background()
	// guest
	og, secret, _ := svc.CreateOrderGuest(ctx, OrderCreateRequest{Currency: "NGN", Items: []OrderItemCreate{{MenuItemID: "item-1", Quantity: 1}}})
	// auth
	ao, _ := svc.CreateOrderCustomer(ctx, "cust-1", OrderCreateRequest{Currency: "NGN", Items: []OrderItemCreate{{MenuItemID: "item-3", Quantity: 1}}})
	// access via secret
	_, err := svc.GetOrder(ctx, og.ID, AuthContext{OrderSecret: &secret})
	if err != nil { t.Fatalf("secret access failed: %v", err) }
	// access via wrong secret
	_, err = svc.GetOrder(ctx, og.ID, AuthContext{OrderSecret: ptr("wrong")})
	if err == nil { t.Fatalf("expected forbidden for wrong secret") }
	// access via customer
	_, err = svc.GetOrder(ctx, ao.ID, AuthContext{CustomerID: ptr("cust-1")})
	if err != nil { t.Fatalf("customer access failed: %v", err) }
}

func TestListOrdersByCustomer(t *testing.T) {
	svc := newOrderSvc()
	ctx := context.Background()
	_, _ = svc.CreateOrderCustomer(ctx, "cust-1", OrderCreateRequest{Currency: "NGN", Items: []OrderItemCreate{{MenuItemID: "item-2", Quantity: 1}}})
	_, _ = svc.CreateOrderCustomer(ctx, "cust-1", OrderCreateRequest{Currency: "NGN", Items: []OrderItemCreate{{MenuItemID: "item-4", Quantity: 1}}})
	list, pag, err := svc.ListMyOrders(ctx, "cust-1", 1, 10)
	if err != nil { t.Fatalf("err: %v", err) }
	if len(list) != 2 || pag.Total < 2 { t.Fatalf("expected 2 orders, got %d", len(list)) }
}

func ptr[T any](v T) *T { return &v }
