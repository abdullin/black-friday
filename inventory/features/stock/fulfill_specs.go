package stock

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {
	Define(&Spec{
		Name: "can't fulfill non-existent reservation",
		When: &FulfillReq{
			Reservation: u(1),
		},
		ThenError: ErrReservationNotFound,
	})

	Define(&Spec{
		Name: "simple fulfillment",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(2), Name: "Warehouse"},
			&InventoryUpdated{Location: u(2), Product: u(1), OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: u(3),
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 7, Location: u(2)}},
			},
		},
		When: &FulfillReq{
			Reservation: u(3),
			Items:       []*FulfillReq_Item{{Product: u(1), Quantity: 7, Location: u(2)}},
		},
		ThenResponse: &FulfillResp{},
		ThenEvents: []proto.Message{&Fulfilled{
			Reservation: u(3),
			Items:       []*Fulfilled_Item{{Product: u(1), Location: u(2), Removed: 7, OnHand: 3}},
		}},
	})

	Define(&Spec{
		Name: "fulfill reservation in a shipment box",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(2), Name: "Warehouse"},
			&LocationAdded{Uid: u(3), Name: "Shelf", Parent: u(2)},
			&InventoryUpdated{Location: u(3), Product: u(1), OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: u(4),
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 7, Location: u(2)}},
			},
		},
		When: &FulfillReq{
			Reservation: u(4),
			Items:       []*FulfillReq_Item{{Product: u(1), Quantity: 7, Location: u(3)}},
		},
		ThenResponse: &FulfillResp{},
		ThenEvents: []proto.Message{&Fulfilled{
			Reservation: u(4),
			Items:       []*Fulfilled_Item{{Product: u(1), Location: u(3), Removed: 7, OnHand: 3}},
		}},
	})

	Define(&Spec{
		Name: "can't fulfill in a way that breaks availability",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(2), Name: "Warehouse"},
			&LocationAdded{Uid: u(3), Name: "Shelf 1", Parent: u(2)},
			&LocationAdded{Uid: u(4), Name: "Shelf 2", Parent: u(2)},
			&InventoryUpdated{Location: u(3), Product: u(1), OnHandChange: 2, OnHand: 2},
			&InventoryUpdated{Location: u(4), Product: u(1), OnHandChange: 1, OnHand: 1},
			&Reserved{
				Reservation: u(5),
				Code:        "whs",
				Items:       []*Reserved_Item{{Product: u(1), Location: u(2), Quantity: 2}},
			},
			&Reserved{
				Reservation: u(6),
				Code:        "shelf",
				Items:       []*Reserved_Item{{Product: u(1), Location: u(3), Quantity: 1}},
			},
		},
		When: &FulfillReq{
			Reservation: u(5),
			Items:       []*FulfillReq_Item{{Product: u(1), Quantity: 2, Location: u(3)}},
		},
		ThenError: ErrNotEnough,
	})

	Define(&Spec{
		Name: "can fulfill in a way that keeps availability",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(2), Name: "Warehouse"},
			&LocationAdded{Uid: u(3), Name: "Shelf 1", Parent: u(2)},
			&LocationAdded{Uid: u(4), Name: "Shelf 2", Parent: u(2)},
			&InventoryUpdated{Location: u(3), Product: u(1), OnHandChange: 2, OnHand: 2},
			&InventoryUpdated{Location: u(4), Product: u(1), OnHandChange: 1, OnHand: 1},
			&Reserved{
				Reservation: u(5),
				Code:        "whs",
				Items:       []*Reserved_Item{{Product: u(1), Location: u(2), Quantity: 2}},
			},
			&Reserved{
				Reservation: u(6),
				Code:        "shelf",
				Items:       []*Reserved_Item{{Product: u(1), Location: u(3), Quantity: 1}},
			},
		},
		When: &FulfillReq{
			Reservation: u(5),
			Items: []*FulfillReq_Item{
				{Product: u(1), Quantity: 1, Location: u(3)},
				{Product: u(1), Quantity: 1, Location: u(4)},
			},
		},
		ThenEvents: []proto.Message{&Fulfilled{
			Reservation: u(5),
			Items: []*Fulfilled_Item{
				{Product: u(1), Location: u(3), Removed: 1, OnHand: 1},
				{Product: u(1), Location: u(4), Removed: 1, OnHand: 0},
			},
		}},
		ThenResponse: &FulfillResp{},
	})
}
