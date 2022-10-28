package specs

import (
	"context"
	"database/sql"
	"github.com/abdullin/go-seq"
)

type Env struct {
	ctx context.Context
	db  *sql.DB
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
