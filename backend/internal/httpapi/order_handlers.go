package httpapi

import (
	"net/http"
)

// handleCreateOrder allows guest access and returns 201 with empty order and orderSecret.
func handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	// Build empty order per schema with zero values and empty items
	empty := OrderDTO{Items: []OrderItemDTO{}}
	resp := map[string]interface{}{
		"order":       empty,
		"orderSecret": "",
	}
	writeJSON(w, http.StatusCreated, resp)
}

// handleListOrders requires auth (enforced by router middleware) and returns empty list with pagination.
func handleListOrders(w http.ResponseWriter, r *http.Request) {
	resp := OrdersResponse{
		Orders: []OrderDTO{},
		Pagination: PaginationDTO{Page: 1, PageSize: 0, Total: 0},
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleGetOrder allows JWT or X-Order-Secret and returns empty order.
func handleGetOrder(w http.ResponseWriter, r *http.Request) {
	empty := OrderResponse{Items: []OrderItemDTO{}}
	writeJSON(w, http.StatusOK, empty)
}
