package postgres

import (
	"context"
	"strconv"

	"local.dev/foodapp/internal/ports/repositories"
)

type MenuItemRepositoryPG struct{ db *DB }

func NewMenuItemRepositoryPG(db *DB) *MenuItemRepositoryPG { return &MenuItemRepositoryPG{db: db} }

func (r *MenuItemRepositoryPG) ListAvailable(ctx context.Context, categoryID *string, page, pageSize int) ([]repositories.MenuItemRow, int, error) {
	q := `SELECT id, name, description, category_id, price_cents, currency, is_available FROM menu_items WHERE is_available=true`
	args := []any{}
	if categoryID != nil && *categoryID != "" { q += ` AND category_id=$1`; args = append(args, *categoryID) }
	q += ` ORDER BY name`
	// total
	qt := `SELECT count(1) FROM (`+q+`) t`
	var total int
	if err := r.db.QueryRow(ctx, qt, args...).Scan(&total); err != nil { return nil, 0, err }
	if page <= 0 { page = 1 }
	if pageSize < 0 { pageSize = 0 }
	if pageSize > 0 {
		offset := (page-1)*pageSize
		q += ` LIMIT `; q += itoa(pageSize); q += ` OFFSET `; q += itoa(offset)
	}
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	var out []repositories.MenuItemRow
	for rows.Next() {
		var rec repositories.MenuItemRow
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.Description, &rec.CategoryID, &rec.PriceCents, &rec.Currency, &rec.IsAvailable); err != nil { return nil, 0, err }
		out = append(out, rec)
	}
	return out, total, rows.Err()
}

func itoa(n int) string { return strconv.Itoa(n) }
