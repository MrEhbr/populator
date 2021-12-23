package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxConn interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
}

var (
	_ PgxConn = &pgxpool.Pool{}
	_ PgxConn = &pgx.Conn{}
	_ PgxConn = pgx.Tx(nil)
)

func NewPgx(conn PgxConn) *pgxConn {
	return &pgxConn{conn: conn}
}

type pgxConn struct {
	conn PgxConn
}

func (c *pgxConn) BeginFunc(ctx context.Context, f func(Conn) error) (err error) {
	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			err = rollbackErr
		}
	}()

	fErr := f(&pgxConn{tx})
	if fErr != nil {
		_ = tx.Rollback(ctx) // ignore rollback error as there is already an error to return
		return fErr
	}

	return tx.Commit(ctx)
}

func (c *pgxConn) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	res, err := c.conn.Exec(ctx, sql, arguments...)
	if err != nil {
		return 0, err
	}

	return int(res.RowsAffected()), nil
}

func (c *pgxConn) Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (Rows, error) {
	return c.conn.Query(ctx, sql, optionsAndArgs...)
}
