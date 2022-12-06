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
		Name: "fulfill reservation in a shipment box",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(2), Name: "Warehouse"},
			&LocationAdded{Uid: u(3), Name: "Shelf", Parent: u(2)},
			&InventoryUpdated{Location: u(3), Product: u(1), OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: u(4),
				Code:        "sale",
				Items:       []*Stock{{Product: u(1), Quantity: 7, Location: u(2)}},
			},
		},
		When: &FulfillReq{
			Reservation: u(4),
			Items:       []*Stock{{Product: u(1), Quantity: 7, Location: u(3)}},
		},
		ThenResponse: &FulfillResp{},
		ThenEvents: []proto.Message{&Fulfilled{
			Reservation: u(4),
			Items:       []*Stock{{Product: u(1), Quantity: 7, Location: u(3)}},
		}},
	})
}
