package services

import "context"

type OrderService interface {
	ListMyOrders(ctx context.Context, customerID string, page, pageSize int) ([]OrderDTO, Pagination, error)
	CreateOrderGuest(ctx context.Context, req OrderCreateRequest) (OrderDTO, string /*orderSecret*/, error)
	CreateOrderCustomer(ctx context.Context, customerID string, req OrderCreateRequest) (OrderDTO, error)
	GetOrder(ctx context.Context, orderID string, access AuthContext) (OrderDTO, error)
}

type OrderCreateRequest struct {
	Items []OrderItemCreate
	DeliveryAddress Address
	TipCents int
	Currency string
}

type OrderItemCreate struct {
	MenuItemID string
	Quantity int
}

type Address struct {
	Line1 string
	Line2 string
	City string
	State string
	PostalCode string
	Country string
}

type OrderDTO struct{
	ID string
	CustomerID *string
	Status string
	TotalCents int
	Currency string
	DeliveryAddressID *string
	DeliveryInstructions *string
	TipCents int
	DeliveryEtaMinutes *int
	DeliveredAt *string
	CourierName *string
	DeliveryTrackingURL *string
	Items []OrderItemDTO
	CreatedAt string
	UpdatedAt string
}

type OrderItemDTO struct{
	ID string
	OrderID string
	MenuItemID string
	NameSnapshot string
	UnitPriceCents int
	Quantity int
}

type AuthContext struct{
	CustomerID *string
	OrderSecret *string
}
