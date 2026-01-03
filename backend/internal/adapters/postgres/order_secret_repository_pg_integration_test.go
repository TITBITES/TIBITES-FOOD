//go:build postgres_integration

package postgres

import (
	"context"
	"testing"
)

func TestOrderSecretRepositoryPG_StoreAndValidate(t *testing.T) {
	db := connectOrSkip(t)
	defer db.Close()
	applyMigrations(t, db)

	// create dummy order
	_, err := db.Pool.Exec(context.Background(), `INSERT INTO orders (status, total_cents, currency) VALUES ('PLACED', 0, 'NGN')`)
	if err != nil {
		t.Fatalf("insert order: %v", err)
	}
	var orderID string
	if err := db.Pool.QueryRow(context.Background(), `SELECT id FROM orders LIMIT 1`).Scan(&orderID); err != nil {
		t.Fatalf("select order id: %v", err)
	}

	repo := NewOrderSecretRepositoryPG(db)
	if err := repo.StoreOnCreate(context.Background(), orderID, "s1"); err != nil {
		t.Fatalf("store: %v", err)
	}
	ok, err := repo.Validate(context.Background(), orderID, "s1")
	if err != nil || !ok {
		t.Fatalf("validate failed")
	}
}
