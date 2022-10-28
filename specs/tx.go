package specs

import (
	"black-friday/fail"
	"black-friday/inventory/apply"
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
)

// tx provides access to the current transaction context and app in one go
type tx struct {
	ctx    context.Context
	tx     *sql.Tx
	events []proto.Message
}

func (c *tx) QueryHack(q string, args ...any) (*sql.Rows, error) {
	// can we make this prettier?
	return c.tx.QueryContext(c.ctx, q, args...)
}

func (c *tx) GetSeq(name string) int64 {
	return c.LookupInt64("select seq from sqlite_sequence where name=?", name)
}

func (c *tx) LookupInt64(query string, args ...any) int64 {
	row := c.tx.QueryRowContext(c.ctx, query, args...)
	var i int64
	err := row.Scan(&i)
	if err == sql.ErrNoRows {
		return 0
	} else if err != nil {
		panic(fmt.Errorf("sql %s: %w", query, err))
	}

	return i
}

func (c *tx) Rollback() error {
	return c.tx.Rollback()

}
func (c *tx) Commit() error {
	return c.tx.Commit()
}

func (c *tx) Apply(e proto.Message) (error, fail.Code) {

	err := apply.Event(c, e)

	if err != nil {
		extracted, failCode := fail.Extract(err)
		return fmt.Errorf("apply %s: %w", reflect.TypeOf(e).String(), extracted), failCode
	}

	c.events = append(c.events, e)
	return nil, fail.None

}

func (c *tx) Exec(query string, args ...any) error {

	_, err := c.tx.ExecContext(c.ctx, query, args...)

	if err != nil {
		return fmt.Errorf("problem with query '%s': %w", query, err)
	}
	return nil

}
