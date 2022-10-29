package specs

import (
	"black-friday/fx"
	"black-friday/inventory/api"
	"black-friday/inventory/features/locations"
	"black-friday/inventory/features/products"
	"black-friday/inventory/features/stock"
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func dispatch(ctx fx.Tx, m proto.Message) (r proto.Message, err error) {

	switch t := m.(type) {
	case *api.AddLocationsReq:
		r, err = locations.Add(ctx, t)
	case *api.AddProductsReq:
		r, err = products.Add(ctx, t)
	case *api.UpdateInventoryReq:
		r, err = stock.Update(ctx, t)
	case *api.ListLocationsReq:
		r, err = locations.List(ctx, t)
	case *api.GetLocInventoryReq:
		r, err = stock.Query(ctx, t)
	case *api.ReserveReq:
		r, err = stock.Reserve(ctx, t)
	case *api.MoveLocationReq:
		r, err = locations.Move(ctx, t)
	default:
		return nil, fmt.Errorf("missing dispatch for %v", reflect.TypeOf(m))
	}

	if r != nil && reflect.ValueOf(r).IsNil() {
		r = nil
	}
	return r, err
}
