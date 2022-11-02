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
	//	start := time.Now()
	rows, err := c.tx.QueryContext(c.ctx, q, args...)

	return rows, err

}

func (c *tx) GetSeq(name string) int64 {
	var id int64
	c.Scan("select seq from sqlite_sequence where name=?", []any{name}, &id)
	return id
}

func (c *tx) Scan(query string, args []any, dest ...any) bool {
	row := c.tx.QueryRowContext(c.ctx, query, args...)
	err := row.Scan(dest...)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		panic(fmt.Errorf("sql %s: %w", query, err))
	}

	return true
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
