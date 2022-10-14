package inventory

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/proto"
)

type Tx struct {
	tx  *sql.Tx
	ctx context.Context
}

func (c *Tx) Exec(query string, args ...any) error {

	_, err := c.tx.ExecContext(c.ctx, query, args...)

	if err != nil {
		return fmt.Errorf("problem with query '%s': %w", query, err)
	}
	return nil

}

func (c *Tx) QueryUint64(query string, args ...any) (uint64, error) {
	row := c.tx.QueryRowContext(c.ctx, query, args...)
	var i uint64
	err := row.Scan(&i)
	return i, err
}

func (c *Tx) QueryInt64(query string, args ...any) (int64, error) {
	row := c.tx.QueryRowContext(c.ctx, query, args...)
	var i int64
	err := row.Scan(&i)
	return i, err
}

func (s *Tx) Apply(e proto.Message) {
	err := apply(s, e)
	if err != nil {
		// we don't expect to fail. Tests are for that
		panic(fmt.Errorf("Error on event apply: %w", err))

	}
}

func (c *Tx) Rollback() {
	err := c.tx.Rollback()
	if err != nil {
		panic(err)
	}
}

func (c *Tx) Commit() {
	// we don't expect to fail
	err := c.tx.Commit()
	if err != nil {
		panic(err)
	}
}
