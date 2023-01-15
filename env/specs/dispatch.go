package specs

import (
	"black-friday/fx"
	"black-friday/inventory/api"
	"black-friday/inventory/features/locations"
	"black-friday/inventory/features/products"
	"black-friday/inventory/features/stock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func Dispatch(tx fx.Tx, m proto.Message) (r proto.Message, err *status.Status) {

	switch t := m.(type) {
	case *api.AddLocationsReq:
		r, err = locations.Add(tx, t)
	case *api.AddProductsReq:
		r, err = products.Add(tx, t)
	case *api.UpdateInventoryReq:
		r, err = stock.Update(tx, t)
	case *api.ListLocationsReq:
		r, err = locations.List(tx, t)
	case *api.GetLocInventoryReq:
		r, err = stock.Query(tx, t)
	case *api.ReserveReq:
		r, err = stock.Reserve(tx, t)
	case *api.MoveLocationReq:
		r, err = locations.Move(tx, t)
	case *api.FulfillReq:
		r, err = stock.Fulfill(tx, t)
	default:
		return nil, status.Newf(codes.Internal, "missing Dispatch for %v", reflect.TypeOf(m))
	}

	if r != nil && reflect.ValueOf(r).IsNil() {
		r = nil
	}
	return r, err
}
