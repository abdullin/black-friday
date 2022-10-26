package stock

import (
	"black-friday/inventory/api"
	"black-friday/specs"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func init() {
	specs.Add(&specs.S{
		Name: "don't allow negative on-hand",
		Given: []proto.Message{
			&api.LocationAdded{Id: 1, Name: "Shelf"},
			&api.ProductAdded{Id: 1, Sku: "NVidia"},
		},
		When:      &api.UpdateInventoryReq{Product: 1, Location: 1, OnHandChange: -1},
		ThenError: codes.FailedPrecondition,
	})

}
