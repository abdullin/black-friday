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
			&LocationAdded{Uid: u(1), Name: "Shelf"},
			&ProductAdded{Uid: u(2), Sku: "NVidia"},
		},
		When:      &UpdateInventoryReq{Product: u(2), Location: u(1), OnHandChange: -1},
		ThenError: ErrNotEnough,
	})

}
