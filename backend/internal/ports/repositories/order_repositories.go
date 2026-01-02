package repositories

import "context"

type OrderRepository interface {
	Create(ctx context.Context, rec OrderRow, items []OrderItemRow) (OrderRow, error)
	GetByID(ctx context.Context, id string) (OrderRow, []OrderItemRow, error)
	ListByCustomer(ctx context.Context, customerID string, page, pageSize int) ([]OrderRow, int /*total*/, error)
}

type OrderSecretRepository interface {
	StoreOnCreate(ctx context.Context, orderID, orderSecret string) error
	Validate(ctx context.Context, orderID, orderSecret string) (bool, error)
}

type OrderRow struct{
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
	CreatedAt string
	UpdatedAt string
}

type OrderItemRow struct{
	ID string
	OrderID string
	MenuItemID string
	NameSnapshot string
	UnitPriceCents int
	Quantity int
}
