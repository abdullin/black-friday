package node

import (
	"black-friday/fail"
	"black-friday/inventory/apply"
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type tx struct {
	ctx context.Context
	tx  *sql.Tx
}

func (c *tx) GetSeq(name string) int64 {

	var id int64
	c.QueryRow("select seq from sqlite_sequence where name=?", name)(&id)
	return id

}

func (c *tx) Apply(e proto.Message) (error, fail.Code) {
	err := apply.Event(c, e)

	if err != nil {
		extracted, failCode := fail.Extract(err)
		return fmt.Errorf("apply %s: %w", reflect.TypeOf(e).String(), extracted), failCode
	}

	return nil, fail.None
}

func (c *tx) QueryHack(q string, args ...any) (*sql.Rows, error) {

	rows, err := c.tx.QueryContext(c.ctx, q, args...)

	return rows, err
}

func (c *tx) QueryRow(query string, args ...any) func(dest ...any) bool {

	return func(dest ...any) bool {

		row := c.tx.QueryRowContext(c.ctx, query, args...)
		err := row.Scan(dest...)
		if err == sql.ErrNoRows {
			return false
		} else if err != nil {
			panic(fmt.Errorf("sql %s: %w", query, err))
		}

		return true
	}

}

func (c *tx) Exec(query string, args ...any) error {

	_, err := c.tx.ExecContext(c.ctx, query, args...)

	if err != nil {
		return fmt.Errorf("problem with query '%s': %w", query, err)
	}
	return nil
}

func (c *tx) Rollback() error {
	err := c.tx.Rollback()
	if err == sql.ErrTxDone {
		return nil
	}
	return err
}

func (t *tx) Commit() error {
	return t.tx.Commit()
}
