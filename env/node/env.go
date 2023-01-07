package node

import (
	"black-friday/env/tracer"
	"black-friday/fx"
	"black-friday/inventory/db"
	"context"
	"database/sql"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"os"
	"strconv"
	"strings"
)

type Env struct {
	ctx context.Context
	db  *sql.DB

	client *kgo.Client

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

func NewEnv(ctx context.Context, file string) *Env {

	dbs, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Panicln("failed to open DB", err)
	}

	seeds := []string{"159.69.176.118:9092"}
	// One client can both produce and consume!
	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.

	producerId := strconv.FormatInt(int64(os.Getpid()), 10)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.TransactionalID(producerId),
		kgo.DefaultProduceTopic("blackfriday-events"),
	)
	if err != nil {
		panic(fmt.Errorf("Can't create kafka client %w", err))
	}

	err = cl.Ping(ctx)
	if err != nil {
		panic("Failed to ping")
	}

	return &Env{
		ctx:         ctx,
		client:      cl,
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
	if env.client != nil {
		env.client.Close()
	}
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
