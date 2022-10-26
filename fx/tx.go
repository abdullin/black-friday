package fx

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/protobuf/proto"
)

type Tx struct {
	// temporarily public
	Tx     *sql.Tx
	ctx    context.Context
	parent *Tx
	events []proto.Message
}

const NestedTxKey = "tx"

// Stash puts transaction into the context, so that it could be passed
// to dispatch method
func (tx *Tx) Stash(ctx context.Context) context.Context {
	return context.WithValue(ctx, NestedTxKey, tx)
}

func Begin(ctx context.Context, db *sql.DB) *Tx {
	inner, hasParent := ctx.Value(NestedTxKey).(*Tx)

	if hasParent {
		return &Tx{
			Tx:     inner.Tx,
			ctx:    ctx,
			parent: inner,
		}
	}

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		// this is never expected to happen
		panic(fmt.Errorf("failed to create tx: %w", err))
	}
	return &Tx{Tx: tx, ctx: ctx}
}

func (s *Tx) Append(e proto.Message) {

	if s.parent != nil {
		s.parent.events = append(s.parent.events, e)
	} else {
		s.events = append(s.events, e)
	}

}

func (c *Tx) Rollback() {
	if c.parent != nil {
		return
	}
	err := c.Tx.Rollback()
	if err != nil {
		panic(err)
	}
}

func (c *Tx) Commit() {
	if c.parent != nil {
		return
	}
	// we don't expect to fail
	err := c.Tx.Commit()
	if err != nil {
		panic(err)
	}
}

func (c *Tx) TestGetEvents() []proto.Message {
	return c.events
}
func (c *Tx) TestClearEvents() {
	c.events = nil
}
