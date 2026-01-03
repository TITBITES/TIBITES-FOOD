// payment_intent_repository_pg.go
package postgres

import (
	"context"

	"local.dev/foodapp/internal/ports/repositories"
)

type PaymentIntentRepositoryPG struct {
	db *DB
}

func NewPaymentIntentRepositoryPG(db *DB) repositories.PaymentIntentRepository {
	return &PaymentIntentRepositoryPG{db: db}
}

func (r *PaymentIntentRepositoryPG) Create(ctx context.Context, rec repositories.PaymentIntentRow) (repositories.PaymentIntentRow, error) {
	// TODO: implement SQL insert
	return rec, nil
}

func (r *PaymentIntentRepositoryPG) GetByID(ctx context.Context, id string) (repositories.PaymentIntentRow, error) {
	// TODO: implement SQL select
	return repositories.PaymentIntentRow{}, nil
}

func (r *PaymentIntentRepositoryPG) GetByOrder(ctx context.Context, orderID string) (repositories.PaymentIntentRow, error) {
	// TODO: implement SQL select
	return repositories.PaymentIntentRow{}, nil
}
