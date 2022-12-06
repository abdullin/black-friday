package stock

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {
	Define(&Spec{
		Level: 3,
		Name:  "reserve sale with one item",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(2), Name: "Shelf"},
			&InventoryUpdated{Location: u(2), Product: u(1), OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 10},
			},
		},
		ThenResponse: &ReserveResp{Reservation: u(3)},
		ThenEvents: []proto.Message{
			&Reserved{
				Reservation: u(3),
				Code:        "sale",
				Items: []*Reserved_Item{
					{Product: u(1), Quantity: 10},
				},
			},
		},
	})

	Define(&Spec{
		Level: 3,
		Name:  "reserve sale in a specific location",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(2), Name: "Shelf"},
			&InventoryUpdated{Location: u(2), Product: u(1), OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Location:    u(2),
			Items:       []*ReserveReq_Item{{Sku: "GPU", Quantity: 10}},
		},
		ThenResponse: &ReserveResp{Reservation: u(3)},
		ThenEvents: []proto.Message{
			&Reserved{
				Reservation: u(3),
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 10, Location: u(2)}},
			},
		},
	})

	Define(&Spec{
		Level: 3,
		Name:  "reservation codes must be unique",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(1), Name: "Shelf"},
			&InventoryUpdated{Location: u(1), Product: u(1), OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: u(1),
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 1, Location: u(1)}},
			},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Location:    u(1),
			Items:       []*ReserveReq_Item{{Sku: "GPU", Quantity: 1}},
		},
		ThenError: ErrAlreadyExists,
	})

	Define(&Spec{
		Level: 5,
		Name:  "reserve sale in a specific location that doesn't have quantity",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(1), Name: "Shelf"},
			&LocationAdded{Uid: u(2), Name: "Empty"},
			&InventoryUpdated{Location: u(1), Product: u(1), OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Location:    u(2),
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 10},
			},
		},
		ThenError: ErrNotEnough,
	})

	Define(&Spec{
		Level: 1,
		Name:  "reserve non-existent sku",
		When: &ReserveReq{
			Reservation: "test",
			Items: []*ReserveReq_Item{
				{Sku: "sale", Quantity: 1},
			},
		},
		ThenError: ErrProductNotFound,
	})

	Define(&Spec{
		Level: 2,
		Name:  "reserve when onHand isn't enough",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "cola"},
			&LocationAdded{Uid: u(1), Name: "WHS1"},
			&InventoryUpdated{Location: u(1), Product: u(1), OnHandChange: 2, OnHand: 2},
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
		Level: 2,
		Name:  "over-reserve",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "cola"},
			&LocationAdded{Uid: u(1), Name: "WHS1"},
			&InventoryUpdated{Location: u(1), Product: u(1), OnHandChange: 2, OnHand: 2},
			&Reserved{
				Reservation: u(1),
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 1, Location: u(1)}},
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

	Define(&Spec{
		Level: 4,
		Name:  "reserve in a location that contains enough inside",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(2), Name: "Container"},
			&LocationAdded{Uid: u(3), Name: "Box", Parent: u(2)},
			&InventoryUpdated{Location: u(3), Product: u(1), OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Location:    u(2),
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 10},
			},
		},
		ThenResponse: &ReserveResp{Reservation: u(4)},
		ThenEvents: []proto.Message{
			&Reserved{
				Reservation: u(4),
				Code:        "sale",
				Items: []*Reserved_Item{
					{Product: u(1), Quantity: 10, Location: u(2)},
				},
			},
		},
	})

	Define(&Spec{
		Level: 5,
		Name:  "reserve in a location that doesn't have enough inside",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(1), Name: "Container"},
			&LocationAdded{Uid: u(2), Name: "Box", Parent: u(1)},
			&InventoryUpdated{Location: u(2), Product: u(1), OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Location:    u(1),
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 11},
			},
		},
		ThenError: ErrNotEnough,
	})

	Define(&Spec{
		Level: 5,
		Name:  "reserve box while container has a reservation on top of it (enough)",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(2), Name: "Container"},
			&LocationAdded{Uid: u(3), Name: "Box", Parent: u(2)},
			&InventoryUpdated{Location: u(3), Product: u(1), OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: u(4),
				Code:        "sale0",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 5, Location: u(2)}},
			},
		},
		When: &ReserveReq{
			Reservation: "sale2",
			Location:    u(3),
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 4},
			},
		},
		ThenResponse: &ReserveResp{
			Reservation: u(5),
		},
		ThenEvents: []proto.Message{
			&Reserved{
				Reservation: u(5),
				Code:        "sale2",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 4, Location: u(3)}},
			},
		},
	})

	Define(&Spec{
		Level: 5,
		Name:  "reserve box while container has a reservation on top of it (not enough)",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(1), Name: "Container"},
			&LocationAdded{Uid: u(2), Name: "Box", Parent: u(1)},
			&InventoryUpdated{Location: u(2), Product: u(1), OnHandChange: 3, OnHand: 3},
			&Reserved{
				Reservation: u(1),
				Code:        "sale0",
				Items: []*Reserved_Item{
					{
						Product:  u(1),
						Quantity: 2,
						Location: u(1),
					},
				},
			},
		},
		When: &ReserveReq{
			Reservation: "sale2",
			Location:    u(2),
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 2},
			},
		},
		ThenError: ErrNotEnough,
	})

}
