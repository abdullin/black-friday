package node

import (
	"black-friday/env/tracer"
	"black-friday/fail"
	"black-friday/inventory/apply"
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"
)

type tx struct {
	ctx    context.Context
	tx     *sql.Tx
	events int64
	trace  *tracer.Tracer
	env    *Env
}

func (c *tx) GetSeq(name string) int64 {

	var id int64
	c.QueryRow("select seq from sqlite_sequence where name=?", name)(&id)
	return id

}

var EventCount int64

func (c *tx) Apply(e proto.Message) (error, fail.Code) {

	c.trace.Begin(string(e.ProtoReflect().Descriptor().Name()))

	err := apply.Event(c, e)
	c.trace.End()

	c.events += 1

	if err != nil {
		extracted, failCode := fail.Extract(err)
		return fmt.Errorf("apply %s: %w", reflect.TypeOf(e).String(), extracted), failCode
	}

	return nil, fail.None
}

func (c *tx) QueryHack(q string, args ...any) (*sql.Rows, error) {

	c.trace.Begin(q)

	prep := c.env.GetStmt(q)
	if len(prep.Stmt) != 1 {
		log.Panicln("Query expects only 1 statement")
	}

	stmt := c.tx.StmtContext(c.ctx, prep.Stmt[0])
	rows, err := stmt.QueryContext(c.ctx, args...)

	c.trace.End()

	return rows, err
}

func (c *tx) QueryRow(query string, args ...any) func(dest ...any) bool {

	c.trace.Begin(query)

	return func(dest ...any) bool {

		prep := c.env.GetStmt(query)
		if len(prep.Stmt) != 1 {
			log.Panicln("Query expects only 1 statement")
		}

		stmt := c.tx.StmtContext(c.ctx, prep.Stmt[0])

		row := stmt.QueryRowContext(c.ctx, args...)
		err := row.Scan(dest...)

		c.trace.End()
		if err == sql.ErrNoRows {
			return false
		} else if err != nil {
			panic(fmt.Errorf("sql %q %v: %w", query, args, err))
		}

		return true
	}

}

func (c *tx) Exec(query string, args ...any) error {
	c.trace.Begin(query)

	prep := c.env.GetStmt(query)

	pos := 0
	for i, s := range prep.Stmt {

		stmt := c.tx.StmtContext(c.ctx, s)

		slice := args[pos : pos+prep.Count[i]]

		_, err := stmt.ExecContext(c.ctx, slice...)

		pos += prep.Count[i]

		if err != nil {

			c.trace.End()
			return fmt.Errorf("problem with query '%s' (%v): %w", query, args, err)
		}
	}

	c.trace.End()

	return nil
}

func (c *tx) Rollback() error {
	c.Trace().Begin("Rollback")
	err := c.tx.Rollback()
	c.Trace().End()
	if err == sql.ErrTxDone {
		return nil
	}

	return err
}

func (t *tx) Commit() error {
	EventCount += t.events
	t.Trace().Begin("Commit")
	err := t.tx.Commit()

	t.Trace().End()

	return err
}

func (t *tx) Trace() *tracer.Tracer {
	return t.trace
}
