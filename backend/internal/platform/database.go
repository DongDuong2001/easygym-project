package platform

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OpenDB connects and pings. For Supabase's transaction pooler (port 6543) add
// default_query_exec_mode=exec to DATABASE_URL: the pooler cannot keep prepared statements.
func OpenDB(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

// SetStatus moves the row with this code to status, which must be one of allowed.
// table is a constant the owning module passes (e.g. "orders.orders"), never user input.
func SetStatus(ctx context.Context, db *pgxpool.Pool, table, code, status string, allowed []string) error {
	if !slices.Contains(allowed, status) {
		return ValidationError{Field: "status", Message: "Trạng thái không hợp lệ"}
	}
	tag, err := db.Exec(ctx, `update `+table+` set status = $2 where code = $1`, code, status)
	if err == nil && tag.RowsAffected() == 0 {
		err = ErrNotFound
	}
	return err
}

// InTx runs fn in a transaction and rolls back if it returns an error.
func InTx(ctx context.Context, db *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // no-op after Commit
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
