package repositories

import "context"

type CategoryRepository interface {
	ListActive(ctx context.Context) ([]CategoryRow, error)
}

type MenuItemRepository interface {
	ListAvailable(ctx context.Context, categoryID *string, page, pageSize int) ([]MenuItemRow, int /*total*/, error)
}

type CategoryRow struct{
	ID string
	Name string
	Description *string
	SortOrder int
	IsActive bool
}

type MenuItemRow struct{
	ID string
	Name string
	Description string
	CategoryID string
	PriceCents int
	Currency string
	IsAvailable bool
}
