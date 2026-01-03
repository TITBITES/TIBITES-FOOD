package httpapi

import (
	"net/http"
	"strconv"
)

// handleListCategories returns an empty categories list aligned to OpenAPI.
func handleListCategories(w http.ResponseWriter, r *http.Request) {
	resp := CategoriesResponse{Categories: []CategoryDTO{}}
	writeJSON(w, http.StatusOK, resp)
}

// handleListMenuItems returns empty items and zero totals pagination.
func handleListMenuItems(w http.ResponseWriter, r *http.Request) {
	// Parse optional page and pageSize; default page=1, pageSize=0 as placeholder
	page := 1
	pageSize := 0
	if v := r.URL.Query().Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := r.URL.Query().Get("pageSize"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			pageSize = n
		}
	}
	resp := MenuItemsResponse{
		Items: []MenuItemDTO{},
		Pagination: PaginationDTO{Page: page, PageSize: pageSize, Total: 0},
	}
	writeJSON(w, http.StatusOK, resp)
}
