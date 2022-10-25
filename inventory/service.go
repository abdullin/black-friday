package inventory

import (
	. "black-friday/api"
	"black-friday/fx"
	"context"
	"database/sql"
)

type Service struct {
	db *sql.DB
	UnimplementedInventoryServiceServer
}

func (s *Service) GetTx(ctx context.Context) *fx.Tx {
	return fx.Begin(ctx, s.db)
}

func NewService(db *sql.DB) *Service {
	if db == nil {
		panic("db is nil")
	}
	return &Service{
		db: db,
	}
}
