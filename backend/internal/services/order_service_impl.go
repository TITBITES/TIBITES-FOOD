package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"local.dev/foodapp/internal/ports/repositories"
)

type orderService struct{
	orders repositories.OrderRepository
	secrets repositories.OrderSecretRepository
	idSeq int
}

func NewOrderService(orders repositories.OrderRepository, secrets repositories.OrderSecretRepository) OrderService {
	return &orderService{orders: orders, secrets: secrets}
}

func (s *orderService) nextSecret() string {
	s.idSeq++
	return fmt.Sprintf("secret-%d", s.idSeq)
}

func (s *orderService) ListMyOrders(ctx context.Context, customerID string, page, pageSize int) ([]OrderDTO, Pagination, error) {
	rows, total, err := s.orders.ListByCustomer(ctx, customerID, page, pageSize)
	if err != nil { return nil, Pagination{}, err }
	out := make([]OrderDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapOrderRow(r, nil))
	}
	return out, Pagination{Page: page, PageSize: pageSize, Total: total}, nil
}

func (s *orderService) CreateOrderGuest(ctx context.Context, req OrderCreateRequest) (OrderDTO, string, error) {
	row := repositories.OrderRow{
		CustomerID: nil,
		Status:     "PLACED",
		TotalCents: 0,
		Currency:   req.Currency,
		TipCents:   req.TipCents,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	created, err := s.createRow(ctx, row, req)
	if err != nil { return OrderDTO{}, "", err }
	secret := s.nextSecret()
	if err := s.secrets.StoreOnCreate(ctx, created.ID, secret); err != nil { return OrderDTO{}, "", err }
	return mapOrderRow(created, req.Items), secret, nil
}

func (s *orderService) CreateOrderCustomer(ctx context.Context, customerID string, req OrderCreateRequest) (OrderDTO, error) {
	row := repositories.OrderRow{
		CustomerID: &customerID,
		Status:     "PLACED",
		TotalCents: 0,
		Currency:   req.Currency,
		TipCents:   req.TipCents,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	created, err := s.createRow(ctx, row, req)
	if err != nil { return OrderDTO{}, err }
	return mapOrderRow(created, req.Items), nil
}

func (s *orderService) createRow(ctx context.Context, row repositories.OrderRow, req OrderCreateRequest) (repositories.OrderRow, error) {
	// map items into rows with zero pricing as per constraints
	items := make([]repositories.OrderItemRow, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, repositories.OrderItemRow{
			MenuItemID: it.MenuItemID,
			Quantity:   it.Quantity,
			NameSnapshot: "",
			UnitPriceCents: 0,
		})
	}
	return s.orders.Create(ctx, row, items)
}

func (s *orderService) GetOrder(ctx context.Context, orderID string, access AuthContext) (OrderDTO, error) {
	row, items, err := s.orders.GetByID(ctx, orderID)
	if err != nil { return OrderDTO{}, err }
	// Access control: customer ownership or valid order secret
	if access.CustomerID != nil {
		if row.CustomerID == nil || *row.CustomerID != *access.CustomerID {
			return OrderDTO{}, errors.New("forbidden")
		}
	} else if access.OrderSecret != nil {
		ok, err := s.secrets.Validate(ctx, orderID, *access.OrderSecret)
		if err != nil { return OrderDTO{}, err }
		if !ok { return OrderDTO{}, errors.New("forbidden") }
	} else {
		return OrderDTO{}, errors.New("unauthorized")
	}
	// map to DTO
	dto := mapOrderRow(row, nil)
	// attach items
	dto.Items = make([]OrderItemDTO, 0, len(items))
	for _, it := range items {
		dto.Items = append(dto.Items, OrderItemDTO{
			ID: it.ID, OrderID: it.OrderID, MenuItemID: it.MenuItemID, NameSnapshot: it.NameSnapshot, Quantity: it.Quantity, UnitPriceCents: it.UnitPriceCents,
		})
	}
	return dto, nil
}

func mapOrderRow(r repositories.OrderRow, _ []OrderItemCreate) OrderDTO {
	return OrderDTO{
		ID: r.ID,
		CustomerID: r.CustomerID,
		Status: r.Status,
		TotalCents: r.TotalCents,
		Currency: r.Currency,
		DeliveryAddressID: r.DeliveryAddressID,
		DeliveryInstructions: r.DeliveryInstructions,
		TipCents: r.TipCents,
		DeliveryEtaMinutes: r.DeliveryEtaMinutes,
		DeliveredAt: r.DeliveredAt,
		CourierName: r.CourierName,
		DeliveryTrackingURL: r.DeliveryTrackingURL,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Items: []OrderItemDTO{},
	}
}
