package postgres

import (
	"context"

	"local.dev/foodapp/internal/ports/repositories"
)

type OrderRepositoryPG struct {
	db *DB
}

// IMPORTANT: constructor returns the INTERFACE, not the concrete type
func NewOrderRepositoryPG(db *DB) repositories.OrderRepository {
	return &OrderRepositoryPG{db: db}
}

func (r *OrderRepositoryPG) Create(
	ctx context.Context,
	rec repositories.OrderRow,
	items []repositories.OrderItemRow,
) (repositories.OrderRow, error) {
	panic("postgres OrderRepositoryPG.Create not implemented yet")
}

func (r *OrderRepositoryPG) GetByID(
	ctx context.Context,
	id string,
) (repositories.OrderRow, []repositories.OrderItemRow, error) {
	panic("postgres OrderRepositoryPG.GetByID not implemented yet")
}

func (r *OrderRepositoryPG) ListByCustomer(
	ctx context.Context,
	customerID string,
	page, pageSize int,
) ([]repositories.OrderRow, int, error) {
	panic("postgres OrderRepositoryPG.ListByCustomer not implemented yet")
}
