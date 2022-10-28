package specs

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/abdullin/go-seq"
)

type Env struct {
	ctx context.Context
	db  *sql.DB
}

func (env *Env) Begin() (*tx, error) {

	dbtx, err := env.db.BeginTx(env.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	ttx := &tx{
		ctx:    env.ctx,
		tx:     dbtx,
		events: nil,
	}

	return ttx, nil
}

func NewEnv(ctx context.Context, db *sql.DB) *Env {
	return &Env{
		ctx: ctx,
		db:  db,
	}
}

type SpecResult struct {
	EventCount int
	Deltas     seq.Issues
}
