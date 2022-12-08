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

	Define(&Spec{
		Name: "can't store things in root",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "NVidia"},
		},
		When:      &UpdateInventoryReq{Product: u(1), Location: u(0), OnHandChange: 1},
		ThenError: ErrArgument,
	})

	Define(&Spec{
		Name: "add items to a location",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Shelf"},
			&ProductAdded{Uid: u(2), Sku: "NVidia"},
		},
		When:         &UpdateInventoryReq{Location: u(1), Product: u(2), OnHandChange: 7},
		ThenResponse: &UpdateInventoryResp{OnHand: 7},
		ThenEvents: []proto.Message{
			&InventoryUpdated{
				Location:     u(1),
				Product:      u(2),
				OnHandChange: 7,
				OnHand:       7,
			},
		},
	})
	Define(&Spec{
		Name: "add items to a location twice",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Shelf"},
			&ProductAdded{Uid: u(2), Sku: "NVidia"},
			&InventoryUpdated{
				Location:     u(1),
				Product:      u(2),
				OnHandChange: 7,
				OnHand:       7,
			},
		},
		When:         &UpdateInventoryReq{Location: u(1), Product: u(2), OnHandChange: 3},
		ThenResponse: &UpdateInventoryResp{OnHand: 10},
		ThenEvents: []proto.Message{
			&InventoryUpdated{
				Location:     u(1),
				Product:      u(2),
				OnHandChange: 3,
				OnHand:       10,
			},
		},
	})

}
