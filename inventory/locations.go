package inventory

import (
	"context"
	"google.golang.org/protobuf/proto"
	"sdk-go/protos"
)

func re[M proto.Message](m M, err error) (M, error) {
	return m, err
}

func (s *Service) ListLocations(ctx context.Context, req *protos.ListLocationsReq) (r *protos.ListLocationsResp, e error) {

	tx := s.GetTx(ctx)

	rows, err := tx.tx.QueryContext(ctx, "SELECT Id, Name FROM Locations")
	if err != nil {
		return re(r, err)
	}
	defer rows.Close()

	var results []*protos.ListLocationsResp_Loc

	for rows.Next() {
		var id uint64
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return re(r, err)
		}
		results = append(results, &protos.ListLocationsResp_Loc{
			Id:   id,
			Name: name,
		})
	}
	return &protos.ListLocationsResp{Locs: results}, nil
}

func (s *Service) AddLocations(ctx context.Context, req *protos.AddLocationsReq) (r *protos.AddLocationsResp, e error) {

	tx := s.GetTx(ctx)

	id, err := tx.QueryUint64("select seq from sqlite_sequence where name='Locations'")

	if err != nil {
		return re(r, err)
	}

	results := make([]uint64, len(req.Names))
	for i, name := range req.Names {
		id += 1

		e := &protos.LocationAdded{
			Name: name,
			Id:   id,
		}
		results[i] = id

		tx.Apply(e)
	}

	tx.Commit()
	return &protos.AddLocationsResp{Ids: results}, nil
}
