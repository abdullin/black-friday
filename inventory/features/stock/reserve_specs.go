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
	lamdaChannelSwitch := []proto.Message{
		&LocationAdded{Id: 1, Name: "WHS1"},
		&LocationAdded{Id: 2, Name: "WHS2"},
		&ProductAdded{Id: 1, Sku: "GPU"},
		&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 5, OnHand: 5},
		&InventoryUpdated{Location: 2, Product: 1, OnHandChange: 5, OnHand: 5},
		&LambdaInstalled{
			Type: Lambda_RESERVE,
			Code: `
	channel = order.tags.channel
	if channel == 'store' then
		reserveAll("WHS1")
	else
		reserveAll("WHS2")
	end`,
		},
	}

	Define(&Spec{
		Name:  "lambda switch on tags",
		Given: lamdaChannelSwitch,
		When: &ReserveReq{
			Reservation: "sale1",
			Items:       []*ReserveReq_Item{{Sku: "GPU", Quantity: 1}},
			Tags:        map[string]string{"channel": "online"},
		},
		ThenResponse: &ReserveResp{Reservation: 1},
		ThenEvents: []proto.Message{&Reserved{
			Reservation: 1,
			Code:        "sale1",
			Items:       []*Reserved_Item{{Product: 1, Quantity: 1, Location: 2}},
		}},
	})

	Define(&Spec{
		Name:  "lambda switch on tags - reverse",
		Given: lamdaChannelSwitch,
		When: &ReserveReq{
			Reservation: "sale1",
			Items:       []*ReserveReq_Item{{Sku: "GPU", Quantity: 1}},
			Tags:        map[string]string{"channel": "store"},
		},
		ThenResponse: &ReserveResp{Reservation: 1},
		ThenEvents: []proto.Message{&Reserved{
			Reservation: 1,
			Code:        "sale1",
			Items:       []*Reserved_Item{{Product: 1, Quantity: 1, Location: 1}},
		}},
	})
}
