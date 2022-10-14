package inventory

import (
	"context"
	"database/sql"
	. "sdk-go/protos"
)

type Loc struct {
	Id   uint64
	Name string
}

type product struct {
	name     string
	quantity map[uint64]int64
}

type Service struct {
	db *sql.DB

	UnimplementedInventoryServiceServer
}

func (s *Service) GetTx(ctx context.Context) *Tx {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		panic("Failed to create tx")
	}
	if tx == nil {
		panic("no tx :(")
	}

	return &Tx{tx: tx, ctx: ctx}

}

func (s *Service) Reset(ctx context.Context, empty *Empty) (*Empty, error) {
	//TODO implement me

	return nil, nil
}

func NewService(db *sql.DB) *Service {
	if db == nil {
		panic("db is nil")
	}
	return &Service{
		db: db,
	}
}
