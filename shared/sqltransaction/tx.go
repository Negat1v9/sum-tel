package sqltransaction

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SqlTx interface {
	StartTx(ctx context.Context) (Txx, error)
}

type Txx interface {
	Commit() error
	Rollback() error
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

type sqlTransaction struct {
	txx *sqlx.Tx
	db  *sqlx.DB
}

func NewSqlTransaction(db *sqlx.DB) SqlTx {
	return &sqlTransaction{
		db: db,
	}
}

func (t *sqlTransaction) StartTx(ctx context.Context) (Txx, error) {
	txx, err := t.db.BeginTxx(ctx, &sql.TxOptions{})
	t.txx = txx
	return t.txx, err
}

// All methods below is adapted from the sqlx.Tx struct
func (t *sqlTransaction) Commit() error {
	return t.txx.Commit()
}

func (t *sqlTransaction) Rollback() error {
	return t.txx.Rollback()
}

func (t *sqlTransaction) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	return t.txx.QueryRowxContext(ctx, query, args...)
}
func (t *sqlTransaction) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return t.txx.GetContext(ctx, dest, query, args...)
}

func (t *sqlTransaction) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.txx.ExecContext(ctx, query, args...)
}

func (t *sqlTransaction) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return t.txx.SelectContext(ctx, dest, query, args...)
}
