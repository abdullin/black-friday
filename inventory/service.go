package inventory

import (
	"context"
	"database/sql"
	"google.golang.org/protobuf/proto"
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

func (s *Service) Apply(tx *sql.Tx, e proto.Message) error {
	apply(tx, e)
	return nil
}

func (s *Service) ApplyEvents(events []proto.Message) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, e := range events {
		apply(tx, e)
	}
	return tx.Commit()
}
