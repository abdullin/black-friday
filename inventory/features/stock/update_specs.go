package stock

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {
	Define(&Spec{
		Level: 2,
		Name:  "don't allow negative on-hand",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Shelf"},
			&ProductAdded{Id: 1, Sku: "NVidia"},
		},
		When:      &UpdateInventoryReq{Product: 1, Location: 1, OnHandChange: -1},
		ThenError: ErrNotEnough,
	})

}
