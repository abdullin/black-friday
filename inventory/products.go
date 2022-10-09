package inventory

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sdk-go/protos"
)

func (s *Service) AddProducts(_ context.Context, req *protos.AddProductsReq) (*protos.AddProductsResp, error) {
	results := make([]uint64, len(req.Skus))
	for i, sku := range req.Skus {
		if _, found := s.store.products_index[sku]; found {
			return nil, status.Errorf(codes.AlreadyExists, "duplicate SKU '%s'", sku)
		}

		id := s.store.prod_counter + 1
		e := &protos.ProductAdded{
			Id:  id,
			Sku: sku,
		}
		s.store.Apply(e)
		results[i] = id
	}
	return &protos.AddProductsResp{Ids: results}, nil
}
