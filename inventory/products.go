package inventory

import (
	"black-friday/api"
	"context"
)

func (s *Service) AddProducts(ctx context.Context, req *api.AddProductsReq) (r *api.AddProductsResp, err error) {

	tx := s.GetTx(ctx)

	id := tx.GetSeq("Products")

	results := make([]uint64, len(req.Skus))
	for i, sku := range req.Skus {

		id += 1
		e := &api.ProductAdded{
			Id:  id,
			Sku: sku,
		}

		s.Apply(tx, e)

		results[i] = id
	}

	tx.Commit()
	return &api.AddProductsResp{Ids: results}, nil
}
