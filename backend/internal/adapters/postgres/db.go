package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct{ Pool *pgxpool.Pool }

func Connect(ctx context.Context, url string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil { return nil, err }
	cfg.MaxConns = 5
	cfg.MinConns = 0
	cfg.MaxConnLifetime = time.Hour
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil { return nil, err }
	if err := pool.Ping(ctx); err != nil { pool.Close(); return nil, err }
	return &DB{Pool: pool}, nil
}

func MustEnvDatabaseURL() string {
	v := os.Getenv("DATABASE_URL")
	if v == "" {
		panic("DATABASE_URL is required in Phase 5.0 runtime")
	}
	return v
}

func (db *DB) Close() { if db.Pool != nil { db.Pool.Close() } }

func (db *DB) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := db.Pool.Exec(ctx, sql, args...)
	return err
}

func (db *DB) Query(ctx context.Context, sql string, args ...any) (pgxRows, error) {
	rows, err := db.Pool.Query(ctx, sql, args...)
	return rows, err
}

func (db *DB) QueryRow(ctx context.Context, sql string, args ...any) pgxRow {
	return db.Pool.QueryRow(ctx, sql, args...)
}

type pgxRow interface { Scan(dest ...any) error }

type pgxRows interface { Next() bool; Scan(dest ...any) error; Err() error; Close() }

func EnsureExtensions(ctx context.Context, db *DB) error {
	// Enable pgcrypto for gen_random_uuid
	if err := db.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS pgcrypto"); err != nil { return fmt.Errorf("enable pgcrypto: %w", err) }
	return nil
}
