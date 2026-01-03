package services

import (
	"context"
	"testing"

	repo "local.dev/foodapp/internal/adapters/repositories/memory"
)

func TestMenuService_ListCategories(t *testing.T) {
	cats := repo.NewCategoryRepoMemory()
	items := repo.NewMenuItemRepoMemory()
	svc := NewMenuService(cats, items)
	ctx := context.Background()
	got, err := svc.ListCategories(ctx)
	if err != nil { t.Fatalf("err: %v", err) }
	if len(got) < 2 { t.Fatalf("expected at least 2 categories, got %d", len(got)) }
}

func TestMenuService_ListMenuItems_All(t *testing.T) {
	cats := repo.NewCategoryRepoMemory()
	items := repo.NewMenuItemRepoMemory()
	svc := NewMenuService(cats, items)
	ctx := context.Background()
	list, pag, err := svc.ListMenuItems(ctx, nil, 0, 0)
	if err != nil { t.Fatalf("err: %v", err) }
	if len(list) < 3 { t.Fatalf("expected some items, got %d", len(list)) }
	if pag.Total < len(list) { t.Fatalf("expected total >= len(list)") }
}

func TestMenuService_ListMenuItems_FilterAndPaginate(t *testing.T) {
	cats := repo.NewCategoryRepoMemory()
	items := repo.NewMenuItemRepoMemory()
	svc := NewMenuService(cats, items)
	ctx := context.Background()
	cat := "cat-1"
	list, pag, err := svc.ListMenuItems(ctx, &cat, 1, 2)
	if err != nil { t.Fatalf("err: %v", err) }
	if len(list) > 2 { t.Fatalf("expected page size <= 2, got %d", len(list)) }
	if pag.Page != 1 || pag.PageSize != 2 { t.Fatalf("unexpected pagination: %+v", pag) }
}
