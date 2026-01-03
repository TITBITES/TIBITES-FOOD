//go:build postgres_integration

package postgres

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"local.dev/foodapp/internal/ports/repositories"
)

func connectOrSkip(t *testing.T) *DB {
	t.Helper()
	url := getenvDBURL()
	if url == "" {
		t.Skip("DATABASE_URL not set; skipping Postgres integration tests")
	}
	db, err := Connect(context.Background(), url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	return db
}

func getenvDBURL() string {
	return os.Getenv("DATABASE_URL")
}

func applyMigrations(t *testing.T, db *DB) {
	t.Helper()
	// apply all .sql files in the migrations directory (sorted)
	dir := "migrations"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read migrations dir: %v", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	for _, name := range names {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("read migration %s: %v", name, err)
		}
		if err := db.Exec(context.Background(), string(b)); err != nil {
			t.Fatalf("apply migration %s: %v", name, err)
		}
	}
}

func TestOrderRepositoryPG_CreateAndGet(t *testing.T) {
	db := connectOrSkip(t)
	defer db.Close()
	applyMigrations(t, db)

	repo := NewOrderRepositoryPG(db)
	row, err := repo.Create(context.Background(), repositories.OrderRow{
		Status:     "PLACED",
		TotalCents: 0,
		Currency:   "NGN",
	}, []repositories.OrderItemRow{
		{MenuItemID: "00000000-0000-0000-0000-000000000001", Quantity: 1},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, items, err := repo.GetByID(context.Background(), row.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID == "" || len(items) != 1 {
		t.Fatalf("unexpected result")
	}
}
