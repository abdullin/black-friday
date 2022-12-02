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
		Name: "reserve sale in a specific location",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "GPU"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Location:    1,
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
		Name: "reserve sale in a specific location that doesn't have quantity",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "GPU"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&LocationAdded{Id: 2, Name: "Empty"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Location:    2,
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 10},
			},
		},
		ThenError: ErrNotEnough,
	})

	Define(&Spec{
		Name: "reserve in a location that contains enough inside",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "GPU"},
			&LocationAdded{Id: 1, Name: "Container"},
			&LocationAdded{Id: 2, Name: "Box", Parent: 1},
			&InventoryUpdated{Location: 2, Product: 1, OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Location:    1,
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 10},
			},
		},
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
		Name: "reserve in a location that doesn't have enough inside",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "GPU"},
			&LocationAdded{Id: 1, Name: "Container"},
			&LocationAdded{Id: 2, Name: "Box", Parent: 1},
			&InventoryUpdated{Location: 2, Product: 1, OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Location:    1,
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 11},
			},
		},
		ThenError: ErrNotEnough,
	})

	Define(&Spec{
		Name: "reserve box while container has a reservation on top of it (enough)",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "GPU"},
			&LocationAdded{Id: 1, Name: "Container"},
			&LocationAdded{Id: 2, Name: "Box", Parent: 1},
			&InventoryUpdated{Location: 2, Product: 1, OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: 1,
				Code:        "sale0",
				Items: []*Reserved_Item{
					{
						Product:  1,
						Quantity: 5,
						Location: 1,
					},
				},
			},
		},
		When: &ReserveReq{
			Reservation: "sale2",
			Location:    2,
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 4},
			},
		},
		ThenResponse: &ReserveResp{
			Reservation: 2,
		},
		ThenEvents: []proto.Message{
			&Reserved{
				Reservation: 1,
				Code:        "sale2",
				Items: []*Reserved_Item{
					{
						Product:  1,
						Quantity: 4,
						Location: 2,
					},
				},
			},
		},
	})

	Define(&Spec{
		Name: "reserve box while container has a reservation on top of it (not enough)",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "GPU"},
			&LocationAdded{Id: 1, Name: "Container"},
			&LocationAdded{Id: 2, Name: "Box", Parent: 1},
			&InventoryUpdated{Location: 2, Product: 1, OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: 1,
				Code:        "sale0",
				Items: []*Reserved_Item{
					{
						Product:  1,
						Quantity: 5,
						Location: 1,
					},
				},
			},
		},
		When: &ReserveReq{
			Reservation: "sale2",
			Location:    2,
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 6},
			},
		},
		ThenError: ErrNotEnough,
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
