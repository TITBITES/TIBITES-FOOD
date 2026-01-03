package memory

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"

	"local.dev/foodapp/internal/ports/repositories"
)

type OrderRepoMemory struct{
	mu    sync.RWMutex
	seq   int
	rows  map[string]repositories.OrderRow
	items map[string][]repositories.OrderItemRow
}

func NewOrderRepoMemory() *OrderRepoMemory {
	return &OrderRepoMemory{rows: map[string]repositories.OrderRow{}, items: map[string][]repositories.OrderItemRow{}}
}

func (m *OrderRepoMemory) nextID() string {
	m.seq++
	return "order-" + strconv.Itoa(m.seq)
}

func (m *OrderRepoMemory) Create(ctx context.Context, rec repositories.OrderRow, items []repositories.OrderItemRow) (repositories.OrderRow, error) {
	m.mu.Lock(); defer m.mu.Unlock()
	if rec.ID == "" { rec.ID = m.nextID() }
	if rec.CreatedAt == "" { rec.CreatedAt = time.Now().UTC().Format(time.RFC3339) }
	rec.UpdatedAt = rec.CreatedAt
	for i := range items {
		if items[i].ID == "" { items[i].ID = rec.ID + "-item-" + strconv.Itoa(i+1) }
		items[i].OrderID = rec.ID
	}
	m.rows[rec.ID] = rec
	m.items[rec.ID] = append([]repositories.OrderItemRow{}, items...)
	return rec, nil
}

func (m *OrderRepoMemory) GetByID(ctx context.Context, id string) (repositories.OrderRow, []repositories.OrderItemRow, error) {
	m.mu.RLock(); defer m.mu.RUnlock()
	rec, ok := m.rows[id]
	if !ok { return repositories.OrderRow{}, nil, errNotFound }
	return rec, append([]repositories.OrderItemRow{}, m.items[id]...), nil
}

func (m *OrderRepoMemory) ListByCustomer(ctx context.Context, customerID string, page, pageSize int) ([]repositories.OrderRow, int, error) {
	m.mu.RLock(); defer m.mu.RUnlock()
	var list []repositories.OrderRow
	for _, r := range m.rows {
		if r.CustomerID != nil && *r.CustomerID == customerID {
			list = append(list, r)
		}
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].CreatedAt < list[j].CreatedAt })
	total := len(list)
	if page <= 0 { page = 1 }
	if pageSize <= 0 { return list, total, nil }
	start := (page-1)*pageSize
	if start >= len(list) { return []repositories.OrderRow{}, total, nil }
	end := start + pageSize
	if end > len(list) { end = len(list) }
	return list[start:end], total, nil
}

var errNotFound = errors.New("not found")
