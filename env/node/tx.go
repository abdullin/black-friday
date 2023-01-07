package node

import (
	"black-friday/env/tracer"
	"black-friday/fail"
	"black-friday/inventory/apply"
	"black-friday/inventory/features/graphs"
	"context"
	"database/sql"
	"fmt"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"
)

type tx struct {
	ctx    context.Context
	tx     *sql.Tx
	events []proto.Message
	trace  *tracer.Tracer
	env    *Env
}

func (c *tx) GetSeq(name string) int64 {

	var id int64
	c.QueryRow("select seq from sqlite_sequence where name=?", name)(&id)
	return id

}

var EventCount int

func (c *tx) Apply(e proto.Message, batch bool) (error, fail.Code) {

	c.trace.Begin(string(e.ProtoReflect().Descriptor().Name()))

	err := apply.Event(c, e, batch)
	c.trace.End()

	c.events = append(c.events, e)

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
	if len(t.events) == 0 {
		panic("Nothing to commit!")

	}

	EventCount += len(t.events)

	t.Trace().Begin("Publish")
	records := make([]*kgo.Record, 0, len(t.events))

	byteCount := 0

	for _, e := range t.events {

		text := prototext.Format(e)
		header := string(e.ProtoReflect().Descriptor().Name())
		full := fmt.Sprintf("%s %s", header, text)

		record := kgo.KeyStringRecord("tenant-1", full)
		records = append(records, record)

	}
	t.Trace().Arg(map[string]interface{}{"events": len(t.events), "bytes": byteCount})

	if err := t.env.client.BeginTransaction(); err != nil {
		panic(fmt.Errorf("error beginning transaction: %v\n", err))
	}

	results := t.env.client.ProduceSync(t.ctx, records...)
	if results.FirstErr() != nil {
		log.Panicln("Problem publishing", results.FirstErr())
	}

	// Attempt to commit the transaction and explicitly abort if the
	// commit was not attempted.
	switch err := t.env.client.EndTransaction(t.ctx, kgo.TryCommit); err {
	case nil:
	case kerr.OperationNotAttempted:
		panic("rollback")
	default:
		panic(fmt.Errorf("error committing transaction: %v\n", err))
	}
	t.Trace().End()

	t.Trace().Begin("Update model")
	t.Trace().Arg(map[string]interface{}{"events": len(t.events)})
	for _, e := range t.events {
		graphs.World.Apply(e)
	}
	t.Trace().End()

	t.Trace().Begin("Commit Tx")
	err := t.tx.Commit()

	t.Trace().End()

	return err
}

func (t *tx) Trace() *tracer.Tracer {
	return t.trace
}
