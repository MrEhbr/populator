package postgres

import (
	"context"
	"database/sql"
	"errors"
)

var (
	_ SQLConn = &sql.DB{}
	_ Conn    = &txWrapper{}
	_ Rows    = &rowsWrapper{}
)

type SQLConn interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	ExecContext(ctx context.Context, sql string, arguments ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, sql string, optionsAndArgs ...interface{}) (*sql.Rows, error)
}

type rowsWrapper struct {
	rows *sql.Rows
}

func (r *rowsWrapper) Close()                         { _ = r.rows.Close() }
func (r *rowsWrapper) Err() error                     { return r.rows.Err() }
func (r *rowsWrapper) Next() bool                     { return r.rows.Next() }
func (r *rowsWrapper) Scan(dest ...interface{}) error { return r.rows.Scan(dest) }

func NewSQL(conn SQLConn) *dbSQL {
	return &dbSQL{conn: conn}
}

type dbSQL struct {
	conn SQLConn
}

func (c *dbSQL) BeginFunc(ctx context.Context, f func(Conn) error) (err error) {
	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = rollbackErr
		}
	}()

	fErr := f(&txWrapper{tx})
	if fErr != nil {
		_ = tx.Rollback() // ignore rollback error as there is already an error to return
		return fErr
	}

	return tx.Commit()
}

func (c *dbSQL) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	res, err := c.conn.ExecContext(ctx, sql, arguments...)
	if err != nil {
		return 0, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(n), nil
}

func (c *dbSQL) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	rows, err := c.conn.QueryContext(ctx, sql, args...) // nolint: sqlclosecheck
	if err != nil {
		return nil, err
	}

	return &rowsWrapper{rows}, nil
}

type txWrapper struct {
	tx *sql.Tx
}

func (w *txWrapper) BeginFunc(ctx context.Context, f func(Conn) error) (err error) {
	return errors.New("not supported")
}

func (w *txWrapper) Exec(ctx context.Context, sql string, arguments ...interface{}) (int, error) {
	res, err := w.tx.ExecContext(ctx, sql, arguments...)
	if err != nil {
		return 0, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(n), nil
}

func (w *txWrapper) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	rows, err := w.tx.QueryContext(ctx, sql, args...) // nolint: sqlclosecheck
	if err != nil {
		return nil, err
	}

	return &rowsWrapper{rows}, nil
}
