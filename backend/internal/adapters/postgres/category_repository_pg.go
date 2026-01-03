package postgres

import (
	"context"

	"local.dev/foodapp/internal/ports/repositories"
)

type CategoryRepositoryPG struct{ db *DB }

func NewCategoryRepositoryPG(db *DB) *CategoryRepositoryPG { return &CategoryRepositoryPG{db: db} }

func (r *CategoryRepositoryPG) ListActive(ctx context.Context) ([]repositories.CategoryRow, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, description, sort_order, is_active FROM categories WHERE is_active=true ORDER BY sort_order, name`)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []repositories.CategoryRow
	for rows.Next() {
		var rec repositories.CategoryRow
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.Description, &rec.SortOrder, &rec.IsActive); err != nil { return nil, err }
		out = append(out, rec)
	}
	return out, rows.Err()
}
