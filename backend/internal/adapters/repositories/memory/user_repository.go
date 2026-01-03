package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"local.dev/foodapp/internal/ports/repositories"
)

type UserRepoMemory struct{
	mu sync.RWMutex
	byID map[string]repositories.UserRecord
	byEmail map[string]string // email -> id
	seq int
}

func NewUserRepoMemory() *UserRepoMemory {
	return &UserRepoMemory{byID: map[string]repositories.UserRecord{}, byEmail: map[string]string{}, seq: 0}
}

func (m *UserRepoMemory) nextID() string {
	m.seq++
	return "user-" + fmtInt(m.seq)
}

func fmtInt(n int) string { return fmt.Sprintf("%d", n) }

func (m *UserRepoMemory) GetByEmail(ctx context.Context, email string) (repositories.UserRecord, error) {
	m.mu.RLock(); defer m.mu.RUnlock()
	id, ok := m.byEmail[email]
	if !ok { return repositories.UserRecord{}, errors.New("not found") }
	return m.byID[id], nil
}

func (m *UserRepoMemory) Create(ctx context.Context, rec repositories.UserRecord) (repositories.UserRecord, error) {
	m.mu.Lock(); defer m.mu.Unlock()
	if _, exists := m.byEmail[rec.Email]; exists {
		return repositories.UserRecord{}, errors.New("duplicate email")
	}
	if rec.ID == "" { rec.ID = m.nextID() }
	m.byID[rec.ID] = rec
	m.byEmail[rec.Email] = rec.ID
	return rec, nil
}

func (m *UserRepoMemory) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	m.mu.RLock(); defer m.mu.RUnlock()
	_, ok := m.byEmail[email]
	return ok, nil
}
