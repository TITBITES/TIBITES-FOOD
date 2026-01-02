package services

import "context"

type MenuService interface {
	ListCategories(ctx context.Context) ([]CategoryDTO, error)
	ListMenuItems(ctx context.Context, categoryID *string, page, pageSize int) ([]MenuItemDTO, Pagination, error)
}

type CategoryDTO struct{
	ID string
	Name string
	Description *string
	SortOrder int
	IsActive bool
}

type MenuItemDTO struct{
	ID string
	Name string
	Description string
	CategoryID string
	PriceCents int
	Currency string
	IsAvailable bool
}

type Pagination struct{
	Page int
	PageSize int
	Total int
}
