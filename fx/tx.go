package fx

import (
	"black-friday/env/tracer"
	"black-friday/fail"
	"context"
	"database/sql"
	"google.golang.org/protobuf/proto"
)

type Tx interface {
	GetSeq() int64
	Apply(e proto.Message) (error, fail.Code)
	QueryHack(q string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) func(dest ...any) bool
	Exec(sql string, args ...any) error
	Rollback() error
	Commit() error
	Trace() *tracer.Tracer
}

type Transactor interface {
	Begin(c context.Context) (Tx, error)
}
