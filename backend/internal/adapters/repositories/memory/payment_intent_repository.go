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

type PaymentIntentRepoMemory struct{
	mu   sync.RWMutex
	seq  int
	rows map[string]repositories.PaymentIntentRow
	byOrder map[string]string // orderID -> paymentIntentID (one active)
}

func NewPaymentIntentRepoMemory() *PaymentIntentRepoMemory {
	return &PaymentIntentRepoMemory{rows: map[string]repositories.PaymentIntentRow{}, byOrder: map[string]string{}}
}

func (m *PaymentIntentRepoMemory) nextID() string { m.seq++; return "pi-" + strconv.Itoa(m.seq) }

func (m *PaymentIntentRepoMemory) Create(ctx context.Context, rec repositories.PaymentIntentRow) (repositories.PaymentIntentRow, error) {
	m.mu.Lock(); defer m.mu.Unlock()
	if rec.ID == "" { rec.ID = m.nextID() }
	if rec.CreatedAt == "" { rec.CreatedAt = time.Now().UTC().Format(time.RFC3339) }
	rec.UpdatedAt = rec.CreatedAt
	m.rows[rec.ID] = rec
	m.byOrder[rec.OrderID] = rec.ID
	return rec, nil
}

func (m *PaymentIntentRepoMemory) GetByID(ctx context.Context, id string) (repositories.PaymentIntentRow, error) {
	m.mu.RLock(); defer m.mu.RUnlock()
	rec, ok := m.rows[id]
	if !ok { return repositories.PaymentIntentRow{}, errors.New("not found") }
	return rec, nil
}

func (m *PaymentIntentRepoMemory) GetByOrder(ctx context.Context, orderID string) (repositories.PaymentIntentRow, error) {
	m.mu.RLock(); defer m.mu.RUnlock()
	if id, ok := m.byOrder[orderID]; ok {
		return m.rows[id], nil
	}
	return repositories.PaymentIntentRow{}, errors.New("not found")
}

// Optional: list sorted by createdAt (not required by interface)
func (m *PaymentIntentRepoMemory) listSorted() []repositories.PaymentIntentRow {
	m.mu.RLock(); defer m.mu.RUnlock()
	out := make([]repositories.PaymentIntentRow, 0, len(m.rows))
	for _, r := range m.rows { out = append(out, r) }
	sort.SliceStable(out, func(i, j int) bool { return out[i].CreatedAt < out[j].CreatedAt })
	return out
}
