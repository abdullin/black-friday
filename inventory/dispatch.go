package inventory

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
	"sdk-go/protos"
)

func (s *Service) Dispatch(ctx context.Context, m proto.Message) (proto.Message, error) {
	m, err := s.dispatchInner(ctx, m)

	// because m is never nil here, even if the returned value was nil
	if m != nil && reflect.ValueOf(m).IsNil() {
		return nil, err
	}
	return m, err

}

func (s *Service) dispatchInner(ctx context.Context, m proto.Message) (proto.Message, error) {
	switch t := m.(type) {
	case *protos.AddLocationsReq:
		return s.AddLocations(ctx, t)
	case *protos.AddProductsReq:
		return s.AddProducts(ctx, t)
	case *protos.UpdateInventoryReq:
		return s.UpdateInventory(ctx, t)
	case *protos.ListLocationsReq:
		return s.ListLocations(ctx, t)
	case *protos.GetInventoryReq:
		return s.GetInventory(ctx, t)
	default:
		return nil, fmt.Errorf("missing dispatch for %v", reflect.TypeOf(m))
	}
}
