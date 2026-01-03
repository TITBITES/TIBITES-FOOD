package memory

import (
	"context"
	"sync"
)

type OrderSecretRepoMemory struct{
	mu sync.RWMutex
	byOrder map[string]string
	bySecret map[string]string
}

func NewOrderSecretRepoMemory() *OrderSecretRepoMemory {
	return &OrderSecretRepoMemory{byOrder: map[string]string{}, bySecret: map[string]string{}}
}

func (m *OrderSecretRepoMemory) StoreOnCreate(ctx context.Context, orderID, orderSecret string) error {
	m.mu.Lock(); defer m.mu.Unlock()
	m.byOrder[orderID] = orderSecret
	m.bySecret[orderSecret] = orderID
	return nil
}

func (m *OrderSecretRepoMemory) Validate(ctx context.Context, orderID, orderSecret string) (bool, error) {
	m.mu.RLock(); defer m.mu.RUnlock()
	stored, ok := m.byOrder[orderID]
	if !ok { return false, nil }
	return stored == orderSecret, nil
}

func (m *OrderSecretRepoMemory) GetOrderIDBySecret(secret string) (string, bool) {
	m.mu.RLock(); defer m.mu.RUnlock()
	id, ok := m.bySecret[secret]
	return id, ok
}
