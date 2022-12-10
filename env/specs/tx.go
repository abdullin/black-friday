package specs

import (
	"black-friday/env/tracer"
	"black-friday/fail"
	"black-friday/inventory/apply"
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
)

// Tx provides access to the current transaction context and app in one go
type Tx struct {
	ctx    context.Context
	tx     *sql.Tx
	Events []proto.Message
}

func (c *Tx) QueryHack(q string, args ...any) (*sql.Rows, error) {
	// can we make this prettier?
	//	start := time.Now()
	rows, err := c.tx.QueryContext(c.ctx, q, args...)

	return rows, err

}

func (c *Tx) GetSeq(name string) int64 {
	// this is TEST environment, so we assign globally incrementing values
	var id int64
	c.QueryRow("select MAX(seq) from sqlite_sequence")(&id)
	return id

}

func (c *Tx) QueryRow(query string, args ...any) func(dest ...any) bool {

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

func (c *Tx) Rollback() error {

	return c.tx.Rollback()

}
func (c *Tx) Commit() error {
	return c.tx.Commit()
}

func (c *Tx) Apply(e proto.Message) (error, fail.Code) {

	err := apply.Event(c, e)

	if err != nil {
		extracted, failCode := fail.Extract(err)
		return fmt.Errorf("apply %s: %w", reflect.TypeOf(e).String(), extracted), failCode
	}

	c.Events = append(c.Events, e)
	return nil, fail.None

}

func (c *Tx) Trace() *tracer.Tracer {
	return tracer.Disabled
}

func (c *Tx) Exec(query string, args ...any) error {

	_, err := c.tx.ExecContext(c.ctx, query, args...)

	if err != nil {
		return fmt.Errorf("problem with query '%s': %w", query, err)
	}
	return nil

}
