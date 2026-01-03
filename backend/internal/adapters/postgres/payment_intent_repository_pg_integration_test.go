//go:build postgres_integration

package postgres

import (
	"context"
	"testing"

	"local.dev/foodapp/internal/ports/repositories"
)

func TestPaymentIntentRepositoryPG_CreateAndGet(t *testing.T) {
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

	repo := NewPaymentIntentRepositoryPG(db)
	rec, err := repo.Create(context.Background(), repositories.PaymentIntentRow{
		OrderID:          orderID,
		Provider:         "stub",
		ProviderIntentID: "prov-1",
		Status:           "pending",
		AmountCents:      0,
		Currency:         "NGN",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := repo.GetByID(context.Background(), rec.ID)
	if err != nil || got.ID == "" {
		t.Fatalf("get by id failed")
	}
	got2, err := repo.GetByOrder(context.Background(), orderID)
	if err != nil || got2.ID != rec.ID {
		t.Fatalf("get by order failed")
	}
}
