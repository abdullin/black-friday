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

func (s *Service) Reset(ctx context.Context, empty *Empty) (*Empty, error) {
	//TODO implement me

	return nil, nil
}

func NewService(db *sql.DB) *Service {

	return &Service{
		db: db,
	}
}
