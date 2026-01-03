package memory

import (
	"context"
	"sort"

	"local.dev/foodapp/internal/ports/repositories"
)

type CategoryRepoMemory struct{
	rows []repositories.CategoryRow
}

func NewCategoryRepoMemory() *CategoryRepoMemory {
	rows := []repositories.CategoryRow{
		{ID: "cat-1", Name: "Burgers", SortOrder: 1, IsActive: true},
		{ID: "cat-2", Name: "Drinks", SortOrder: 2, IsActive: true},
	}
	return &CategoryRepoMemory{rows: rows}
}

func (m *CategoryRepoMemory) ListActive(ctx context.Context) ([]repositories.CategoryRow, error) {
	out := make([]repositories.CategoryRow, 0, len(m.rows))
	for _, r := range m.rows {
		if r.IsActive { out = append(out, r) }
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].SortOrder == out[j].SortOrder { return out[i].Name < out[j].Name }
		return out[i].SortOrder < out[j].SortOrder
	})
	return out, nil
}
