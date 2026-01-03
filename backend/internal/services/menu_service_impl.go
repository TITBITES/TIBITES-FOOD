package services

import (
	"context"

	"local.dev/foodapp/internal/ports/repositories"
)

type menuService struct{
	cats repositories.CategoryRepository
	items repositories.MenuItemRepository
}

func NewMenuService(cats repositories.CategoryRepository, items repositories.MenuItemRepository) MenuService {
	return &menuService{cats: cats, items: items}
}

func (m *menuService) ListCategories(ctx context.Context) ([]CategoryDTO, error) {
	rows, err := m.cats.ListActive(ctx)
	if err != nil { return nil, err }
	out := make([]CategoryDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, CategoryDTO{
			ID: r.ID,
			Name: r.Name,
			Description: r.Description,
			SortOrder: r.SortOrder,
			IsActive: r.IsActive,
		})
	}
	return out, nil
}

func (m *menuService) ListMenuItems(ctx context.Context, categoryID *string, page, pageSize int) ([]MenuItemDTO, Pagination, error) {
	rows, total, err := m.items.ListAvailable(ctx, categoryID, page, pageSize)
	if err != nil { return nil, Pagination{}, err }
	out := make([]MenuItemDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, MenuItemDTO{
			ID: r.ID,
			Name: r.Name,
			Description: r.Description,
			CategoryID: r.CategoryID,
			PriceCents: r.PriceCents,
			Currency: r.Currency,
			IsAvailable: r.IsAvailable,
		})
	}
	p := Pagination{Page: page, PageSize: pageSize, Total: total}
	return out, p, nil
}
