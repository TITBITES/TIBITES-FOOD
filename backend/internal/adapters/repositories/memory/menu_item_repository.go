package memory

import (
	"context"
	"sort"

	"local.dev/foodapp/internal/ports/repositories"
)

type MenuItemRepoMemory struct{
	rows []repositories.MenuItemRow
}

func NewMenuItemRepoMemory() *MenuItemRepoMemory {
	rows := []repositories.MenuItemRow{
		{ID: "item-1", Name: "Classic Burger", Description: "", CategoryID: "cat-1", PriceCents: 2500, Currency: "NGN", IsAvailable: true},
		{ID: "item-2", Name: "Cheese Burger", Description: "", CategoryID: "cat-1", PriceCents: 3000, Currency: "NGN", IsAvailable: true},
		{ID: "item-3", Name: "Veggie Burger", Description: "", CategoryID: "cat-1", PriceCents: 2700, Currency: "NGN", IsAvailable: true},
		{ID: "item-4", Name: "Cola", Description: "", CategoryID: "cat-2", PriceCents: 800, Currency: "NGN", IsAvailable: true},
		{ID: "item-5", Name: "Water", Description: "", CategoryID: "cat-2", PriceCents: 500, Currency: "NGN", IsAvailable: true},
	}
	return &MenuItemRepoMemory{rows: rows}
}

func (m *MenuItemRepoMemory) ListAvailable(ctx context.Context, categoryID *string, page, pageSize int) ([]repositories.MenuItemRow, int, error) {
	// Filter by category if provided
	filtered := make([]repositories.MenuItemRow, 0, len(m.rows))
	for _, r := range m.rows {
		if !r.IsAvailable { continue }
		if categoryID != nil && *categoryID != "" && r.CategoryID != *categoryID { continue }
		filtered = append(filtered, r)
	}
	// Stable ordering by Name
	sort.SliceStable(filtered, func(i, j int) bool { return filtered[i].Name < filtered[j].Name })
	total := len(filtered)
	// Simple pagination (repo-controlled)
	if page <= 0 { page = 1 }
	if pageSize < 0 { pageSize = 0 }
	if pageSize == 0 {
		// return all
		return filtered, total, nil
	}
	start := (page-1)*pageSize
	if start >= len(filtered) {
		return []repositories.MenuItemRow{}, total, nil
	}
	end := start + pageSize
	if end > len(filtered) { end = len(filtered) }
	return filtered[start:end], total, nil
}
