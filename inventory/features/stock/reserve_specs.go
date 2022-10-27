package stock

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {
	Define(&Spec{
		Name: "reserve sale with one item",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "GPU"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Items: []*ReserveReq_Item{
				{
					Sku:      "GPU",
					Quantity: 10,
				},
			},
		},
		ThenResponse: &ReserveResp{
			Reservation: 1,
		},
		ThenEvents: []proto.Message{
			&Reserved{
				Reservation: 1,
				Code:        "sale",
				Items: []*Reserved_Item{
					{
						Product:  1,
						Quantity: 10,
					},
				},
			},
		},
	})

}
