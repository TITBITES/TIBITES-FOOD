package postgres

import (
	"context"
	"strings"

	"local.dev/foodapp/internal/ports/repositories"
)

type UserRepositoryPG struct{ db *DB }

func NewUserRepositoryPG(db *DB) *UserRepositoryPG { return &UserRepositoryPG{db: db} }

func (r *UserRepositoryPG) GetByEmail(ctx context.Context, email string) (repositories.UserRecord, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	row := r.db.QueryRow(ctx, `SELECT id, email, password_hash, role FROM users WHERE email=$1`, email)
	var rec repositories.UserRecord
	if err := row.Scan(&rec.ID, &rec.Email, &rec.PasswordHash, &rec.Role); err != nil { return repositories.UserRecord{}, err }
	return rec, nil
}

func (r *UserRepositoryPG) Create(ctx context.Context, rec repositories.UserRecord) (repositories.UserRecord, error) {
	rec.Email = strings.ToLower(strings.TrimSpace(rec.Email))
	row := r.db.QueryRow(ctx, `INSERT INTO users (email, password_hash, role) VALUES ($1,$2,$3) RETURNING id, email, password_hash, role`, rec.Email, rec.PasswordHash, rec.Role)
	if err := row.Scan(&rec.ID, &rec.Email, &rec.PasswordHash, &rec.Role); err != nil { return repositories.UserRecord{}, err }
	return rec, nil
}

func (r *UserRepositoryPG) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	row := r.db.QueryRow(ctx, `SELECT 1 FROM users WHERE email=$1`, email)
	var one int
	if err := row.Scan(&one); err != nil { return false, nil }
	return true, nil
}
