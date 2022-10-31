package specs

import (
	"black-friday/inventory/db"
	"context"
	"database/sql"
	"fmt"
	"github.com/abdullin/go-seq"
	"log"
)

type Env struct {
	ctx context.Context
	db  *sql.DB

	schemaReady bool
}

func NewEnv(ctx context.Context, file string) *Env {

	dbs, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Panicln("failed to open DB", err)
	}

	return &Env{
		ctx:         ctx,
		db:          dbs,
		schemaReady: false,
	}
}

func (env *Env) Close() {
	if env.db != nil {
		err := env.db.Close()
		if err != nil {
			log.Panicln("Failed to close db", err)
		}
		env.db = nil
	}
}

func (env *Env) EnsureSchema() {
	if env.schemaReady {
		return
	}
	err := db.CreateSchema(env.db)
	if err != nil {
		log.Panicln("can't prepare schema: ", err)

	}

	env.schemaReady = true

}

func (env *Env) BeginTx() (*tx, error) {

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

type SpecResult struct {
	EventCount int
	Deltas     seq.Issues
}

func (s *SpecResult) DidFail() bool {
	return len(s.Deltas) > 0

}
