package inventory

import (
	"black-friday/api"
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
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
	case *api.AddLocationsReq:
		return s.AddLocations(ctx, t)
	case *api.AddProductsReq:
		return s.AddProducts(ctx, t)
	case *api.UpdateInventoryReq:
		return s.UpdateInventory(ctx, t)
	case *api.ListLocationsReq:
		return s.ListLocations(ctx, t)
	case *api.GetLocInventoryReq:
		return s.GetLocInventory(ctx, t)
	case *api.MoveLocationReq:
		return s.MoveLocation(ctx, t)
	default:
		return nil, fmt.Errorf("missing dispatch for %v", reflect.TypeOf(m))
	}
}
