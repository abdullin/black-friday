package inventory

import (
	"context"
	"errors"
	"github.com/mattn/go-sqlite3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"sdk-go/protos"
)

func re[M proto.Message](m M, err error) (M, error) {

	if err == nil {
		return m, nil
	}

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		switch sqliteErr.Code {
		case sqlite3.ErrConstraint:
			return m, status.Error(codes.FailedPrecondition, "Constraint violation")
		default:
			return m, status.Errorf(codes.Internal, err.Error())
		}
	}

	return m, err
}

func (s *Service) ListLocations(ctx context.Context, req *protos.ListLocationsReq) (r *protos.ListLocationsResp, e error) {

	tx := s.GetTx(ctx)

	rows, err := tx.tx.QueryContext(ctx, `
SELECT Id, Name, Warehouse FROM Locations
WHERE Warehouse=?
`, req.Warehouse)
	if err != nil {
		return re(r, err)
	}
	defer rows.Close()

	var results []*protos.ListLocationsResp_Loc

	for rows.Next() {
		var id uint64
		var warehouse uint32
		var name string
		err := rows.Scan(&id, &name, &warehouse)
		if err != nil {
			return re(r, err)
		}
		results = append(results, &protos.ListLocationsResp_Loc{
			Warehouse: warehouse,
			Location:  id,
			Name:      name,
		})
	}
	return &protos.ListLocationsResp{Locs: results}, nil
}

func (s *Service) AddLocations(ctx context.Context, req *protos.AddLocationsReq) (r *protos.AddLocationsResp, e error) {

	if req.Warehouse == 0 {
		return nil, status.Error(codes.InvalidArgument, "Warehouse id can't be zero")
	}

	tx := s.GetTx(ctx)

	id, err := tx.QueryUint64("select seq from sqlite_sequence where name='Locations'")

	if err != nil {
		return re(r, err)
	}

	results := make([]uint64, len(req.Names))
	for i, name := range req.Names {
		id += 1

		e := &protos.LocationAdded{
			Name:      name,
			Id:        id,
			Warehouse: req.Warehouse,
		}
		results[i] = id

		err = tx.Apply(e)
		if err != nil {
			return re(r, err)
		}
	}

	tx.Commit()
	return &protos.AddLocationsResp{Ids: results, Warehouse: req.Warehouse}, nil
}
