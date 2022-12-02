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
				{Sku: "GPU", Quantity: 10},
			},
		},
		ThenResponse: &ReserveResp{Reservation: 1},
		ThenEvents: []proto.Message{
			&Reserved{
				Reservation: 1,
				Code:        "sale",
				Items: []*Reserved_Item{
					{Product: 1, Quantity: 10, Location: 1},
				},
			},
		},
	})

	Define(&Spec{
		Name: "reserve non-existent sku",
		When: &ReserveReq{
			Reservation: "test",
			Items: []*ReserveReq_Item{
				{Sku: "sale", Quantity: 1},
			},
		},
		ThenError: ErrProductNotFound,
	})

	Define(&Spec{
		Name: "reserve when onHand isn't enough",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "cola"},
			&LocationAdded{Id: 1, Name: "WHS1"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 2, OnHand: 2},
		},
		When: &ReserveReq{
			Reservation: "test",
			Items: []*ReserveReq_Item{
				{Sku: "cola", Quantity: 3},
			},
		},
		ThenError: ErrNotEnough,
	})

	Define(&Spec{
		Name: "over-reserve",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "cola"},
			&LocationAdded{Id: 1, Name: "WHS1"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 2, OnHand: 2},
			&Reserved{
				Reservation: 1,
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: 1, Quantity: 1, Location: 1}},
			},
		},
		When: &ReserveReq{
			Reservation: "test",
			Items: []*ReserveReq_Item{
				{Sku: "cola", Quantity: 2},
			},
		},
		ThenError: ErrNotEnough,
	})

}
