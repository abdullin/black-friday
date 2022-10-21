package inventory

import (
	. "black-friday/api"
	"context"
	"database/sql"
	"fmt"
)

type Loc struct {
	Id   uint64
	Name string
}

type Service struct {
	db *sql.DB

	UnimplementedInventoryServiceServer
}

func (s *Service) GetTx(ctx context.Context) *Tx {

	inner, hasParent := ctx.Value("tx").(*Tx)

	if hasParent {
		return &Tx{
			tx:     inner.tx,
			ctx:    ctx,
			parent: inner,
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		// this is never expected to happen
		panic(fmt.Errorf("failed to create tx: %w", err))
	}

	return &Tx{tx: tx, ctx: ctx}

}

func NewService(db *sql.DB) *Service {
	if db == nil {
		panic("db is nil")
	}
	return &Service{
		db: db,
	}
}
