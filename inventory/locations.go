package inventory

import (
	"context"
	"sdk-go/protos"
)

func (s *Service) ListLocations(ctx context.Context, req *protos.ListLocationsReq) (*protos.ListLocationsResp, error) {

	results := make([]*protos.ListLocationsResp_Loc, len(s.store.locs))
	for i, l := range s.store.locs {
		results[i] = &protos.ListLocationsResp_Loc{
			Id:   l.Id,
			Name: l.Name,
		}

	}
	return &protos.ListLocationsResp{Locs: results}, nil
}

func (s *Service) AddLocations(_ context.Context, req *protos.AddLocationsReq) (*protos.AddLocationsResp, error) {

	results := make([]uint64, len(req.Names))
	for i, name := range req.Names {
		var id = s.store.loc_counter + 1

		e := &protos.LocationAdded{
			Name: name,
			Id:   id,
		}
		results[i] = id

		s.store.Apply(e)
	}

	return &protos.AddLocationsResp{Ids: results}, nil
}
