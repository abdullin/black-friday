package fx

import (
	"black-friday/env/tracer"
	"black-friday/fail"
	"black-friday/inventory/mem"
	"context"
	"database/sql"
	"google.golang.org/protobuf/proto"
)

type Tx interface {
	GetSeq(table string) int64
	Apply(e proto.Message, batch bool) (error, fail.Code)
	QueryHack(q string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) func(dest ...any) bool
	Exec(sql string, args ...any) error
	Rollback() error
	Commit() error
	Trace() *tracer.Tracer

	GetStockModel(i int32) *mem.ProductStock
	LookupProduct(sku string) (int32, bool)
}

type Transactor interface {
	Begin(c context.Context) (Tx, error)
}
