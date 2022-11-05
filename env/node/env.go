package node

import (
	"black-friday/fx"
	"black-friday/inventory/db"
	"context"
	"database/sql"
	"fmt"
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

func (env *Env) Begin(ctx context.Context) (fx.Tx, error) {

	dbtx, err := env.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	ttx := &tx{
		ctx: env.ctx,
		tx:  dbtx,
	}

	return ttx, nil
}
