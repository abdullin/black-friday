package inventory

import (
	"context"
	"sdk-go/protos"
)

func (s *Service) AddProducts(ctx context.Context, req *protos.AddProductsReq) (r *protos.AddProductsResp, err error) {

	tx := s.GetTx(ctx)

	id := tx.GetSeq("Products")

	results := make([]uint64, len(req.Skus))
	for i, sku := range req.Skus {

		id += 1
		e := &protos.ProductAdded{
			Id:  id,
			Sku: sku,
		}

		tx.Apply(e)

		results[i] = id
	}

	tx.Commit()
	return &protos.AddProductsResp{Ids: results}, nil
}
