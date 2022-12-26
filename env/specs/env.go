package specs

import (
	"black-friday/inventory/db"
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Env struct {
	db *sql.DB

	schemaReady bool
}

func NewEnv(file string) *Env {

	dbs, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Panicln("failed to open DB", err)
	}

	return &Env{
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
	err := db.CreateSchema(env.db, true)
	if err != nil {
		log.Panicln("can't prepare schema: ", err)

	}

	env.schemaReady = true

}

func (env *Env) BeginTx(ctx context.Context) (*Tx, error) {

	dbtx, err := env.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin Tx: %w", err)
	}

	ttx := &Tx{
		ctx:    ctx,
		tx:     dbtx,
		Events: nil,
	}

	return ttx, nil
}

func (s *SpecResult) DidFail() bool {
	return len(s.Deltas) > 0

}
