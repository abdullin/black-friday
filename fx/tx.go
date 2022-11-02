package fx

import (
	"black-friday/fail"
	"context"
	"database/sql"
	"google.golang.org/protobuf/proto"
)

type Tx interface {
	GetSeq(table string) int64
	Apply(e proto.Message) (error, fail.Code)
	QueryHack(q string, args ...any) (*sql.Rows, error)
	Scan(q string, args []any, dest ...any) bool
	Exec(sql string, args ...any) error
	Rollback() error
	Commit() error
}

type Transactor interface {
	Begin(c context.Context) (Tx, error)
}
