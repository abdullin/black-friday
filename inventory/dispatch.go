package inventory

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"sdk-go/protos"
)

// used for command dispatch in tests

func (s Service) Dispatch(ctx context.Context, m proto.Message) (proto.Message, error) {
	switch t := m.(type) {
	case *protos.AddLocationsReq:
		return s.AddLocations(ctx, t)
	case *protos.AddProductsReq:
		return s.AddProducts(ctx, t)
	case *protos.UpdateQtyReq:
		return s.UpdateQty(ctx, t)
	case *protos.ListLocationsReq:
		return s.ListLocations(ctx, t)
	case *protos.GetInventoryReq:
		return s.GetInventory(ctx, t)
	default:
		return nil, fmt.Errorf("missing dispatch for %v", t)
	}
}
