package node

import (
	"black-friday/env/tracer"
	"black-friday/fx"
	"black-friday/inventory/db"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Env struct {
	ctx context.Context
	db  *sql.DB

	store fx.EventStore

	schemaReady bool
	Bank        *tracer.Bank
	prepared    map[string]*Prepared
}

type Prepared struct {
	Count []int
	Stmt  []*sql.Stmt
}

func (e *Env) GetStmt(q string) *Prepared {
	s, found := e.prepared[q]

	if found {
		return s

	}

	s = &Prepared{}

	for _, part := range strings.Split(q, ";") {

		if strings.TrimSpace(part) == "" {
			continue
		}
		x, err := e.db.PrepareContext(e.ctx, part)

		if err != nil {
			log.Panicln(err)
		}

		count := strings.Count(part, "?")

		s.Count = append(s.Count, count)
		s.Stmt = append(s.Stmt, x)
	}

	e.prepared[q] = s

	return s

}

func NewEnv(ctx context.Context, file string, store fx.EventStore) *Env {

	dbs, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Panicln("failed to open DB", err)
	}

	return &Env{
		ctx:         ctx,
		store:       store,
		db:          dbs,
		schemaReady: false,
		Bank:        tracer.NewBank(),
		prepared:    make(map[string]*Prepared),
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
	env.store.Close()
}

func (env *Env) EnsureSchema() {
	if env.schemaReady {
		return
	}
	err := db.CreateSchema(env.db, false)
	if err != nil {
		log.Panicln("can't prepare schema: ", err)

	}

	env.schemaReady = true

}

func (env *Env) DB() *sql.DB {
	return env.db
}

func (env *Env) Begin(ctx context.Context) (fx.Tx, error) {

	trace := env.Bank.Open()
	dbtx, err := env.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	ttx := &tx{
		ctx:   env.ctx,
		tx:    dbtx,
		trace: trace,
		env:   env,
	}

	return ttx, nil
}
