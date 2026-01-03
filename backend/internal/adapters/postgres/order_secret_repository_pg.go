// order_secret_repository_pg.go
package postgres

import (
	"context"

	"local.dev/foodapp/internal/ports/repositories"
)

type OrderSecretRepositoryPG struct {
	db *DB
}

func NewOrderSecretRepositoryPG(db *DB) repositories.OrderSecretRepository {
	return &OrderSecretRepositoryPG{db: db}
}

func (r *OrderSecretRepositoryPG) StoreOnCreate(ctx context.Context, orderID, orderSecret string) error {
	// TODO: implement SQL insert
	return nil
}

func (r *OrderSecretRepositoryPG) Validate(ctx context.Context, orderID, orderSecret string) (bool, error) {
	// TODO: implement SQL select
	return false, nil
}
