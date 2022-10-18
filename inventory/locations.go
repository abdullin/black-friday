package inventory

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	. "sdk-go/protos"
)

func (s *Service) ListLocations(ctx context.Context, req *ListLocationsReq) (r *ListLocationsResp, e error) {

	tx := s.GetTx(ctx)

	rows, err := tx.tx.QueryContext(ctx, `
SELECT Id, Name, Warehouse FROM Locations
WHERE Warehouse=?
`, req.Warehouse)
	if err != nil {
		return re(r, err)
	}
	defer rows.Close()

	var results []*ListLocationsResp_Loc

	for rows.Next() {
		var id uint64
		var warehouse uint32
		var name string
		err := rows.Scan(&id, &name, &warehouse)
		if err != nil {
			return re(r, err)
		}
		results = append(results, &ListLocationsResp_Loc{
			Warehouse: warehouse,
			Location:  id,
			Name:      name,
		})
	}
	return &ListLocationsResp{Locs: results}, nil
}

func (s *Service) CreateWarehouse(ctx context.Context, req *CreateWarehouseReq) (r *CreateWarehouseResp, e error) {
	tx := s.GetTx(ctx)

	id := uint32(tx.GetSeq("Warehouses"))

	results := make([]uint32, len(req.Names))
	for i, name := range req.Names {
		id += 1

		e := &WarehouseCreated{
			Name: name,
			Id:   id,
		}
		results[i] = id

		err := tx.Apply(e)
		if err != nil {
			return re(r, err)
		}
	}

	tx.Commit()
	return &CreateWarehouseResp{Ids: results}, nil
}

func (s *Service) AddLocations(ctx context.Context, req *AddLocationsReq) (r *AddLocationsResp, e error) {

	if req.Warehouse == 0 {
		return nil, status.Error(codes.InvalidArgument, "Warehouse id can't be zero")
	}

	tx := s.GetTx(ctx)

	id := tx.GetSeq("Locations")

	results := make([]uint64, len(req.Names))
	for i, name := range req.Names {
		id += 1

		e := &LocationAdded{
			Name:      name,
			Id:        id,
			Warehouse: req.Warehouse,
		}
		results[i] = id

		err := tx.Apply(e)
		if err != nil {
			return re(r, err)
		}
	}

	tx.Commit()
	return &AddLocationsResp{Ids: results, Warehouse: req.Warehouse}, nil
}
