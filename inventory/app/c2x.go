package app

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/proto"
)

// Context provides access to the current transaction context and app in one go
type Context struct {
	app *App
	ctx context.Context
	tx  *sql.Tx

	events []proto.Message
}

func (c *Context) QueryHack(q string, args ...any) (*sql.Rows, error) {
	// can we make this prettier?
	return c.tx.QueryContext(c.ctx, q, args...)
}

func (c *Context) GetSeq(name string) uint64 {
	id, err := c.QueryUint64("select seq from sqlite_sequence where name=?", name)
	if err != nil {
		panic(fmt.Errorf("failed to get seq for '%s': %w", name, err))
	}
	return id
}

func (c *Context) QueryUint64(query string, args ...any) (uint64, error) {
	row := c.tx.QueryRowContext(c.ctx, query, args...)
	var i uint64
	err := row.Scan(&i)
	return i, err
}

func (c *Context) LookupUint64(query string, args ...any) uint64 {
	row := c.tx.QueryRowContext(c.ctx, query, args...)
	var i uint64
	err := row.Scan(&i)
	if err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		panic(fmt.Errorf("sql %s: %w", query, err))
	}
	return i
}

func (c *Context) QueryInt64(query string, args ...any) (int64, error) {
	row := c.tx.QueryRowContext(c.ctx, query, args...)
	var i int64
	err := row.Scan(&i)

	return i, err
}

func (a *App) Begin(ctx context.Context) (*Context, error) {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("can't begin tx: %w", err)
	}
	return &Context{
		app: a,
		ctx: ctx,
		tx:  tx,
	}, nil
}
func (c *Context) Rollback() error {
	return c.tx.Rollback()

}
func (c *Context) Commit() error {
	return c.tx.Commit()
}

func (c *Context) Exec(query string, args ...any) error {

	_, err := c.tx.ExecContext(c.ctx, query, args...)

	if err != nil {
		return fmt.Errorf("problem with query '%s': %w", query, err)
	}
	return nil

}
