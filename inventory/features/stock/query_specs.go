package stock

import (
	"black-friday/env/uid"
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func u(i int64) string {
	return uid.ToTestString(i)
}

func init() {
	Define(&Spec{
		Level: 3,
		Name:  "query inventory at a specific location",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "Cola"},
			&ProductAdded{Uid: u(2), Sku: "Fanta"},
			&LocationAdded{Uid: u(3), Name: "Shelf", Parent: u(0)},
			&InventoryUpdated{Location: u(3), Product: u(2), OnHandChange: 2, OnHand: 2},
		},
		When: &GetLocInventoryReq{Location: u(3)},
		ThenResponse: &GetLocInventoryResp{
			Items: []*GetLocInventoryResp_Item{{Product: u(2), OnHand: 2, Available: 2}}},
	})

	Define(&Spec{
		Level: 3,
		Name:  "two boxes sum up their quantity at root",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Shelf1", Parent: u(0)},
			&LocationAdded{Uid: u(2), Name: "Shelf2", Parent: u(0)},
			&ProductAdded{Uid: u(3), Sku: "Epyc"},
			&InventoryUpdated{Location: u(1), Product: u(3), OnHandChange: 2, OnHand: 2},
			&InventoryUpdated{Location: u(2), Product: u(3), OnHandChange: 3, OnHand: 3},
		},
		When: &GetLocInventoryReq{Location: u(0)},
		ThenResponse: &GetLocInventoryResp{
			Items: []*GetLocInventoryResp_Item{{Product: u(3), OnHand: 5, Available: 5}}},
	})
	Define(&Spec{
		Level: 3,
		Name:  "boxes sums up quantity with parent container",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Shelf", Parent: u(0)},
			&LocationAdded{Uid: u(2), Name: "Bin", Parent: u(1)},
			&ProductAdded{Uid: u(3), Sku: "Epyc"},
			&InventoryUpdated{Location: u(1), Product: u(3), OnHandChange: 2, OnHand: 2},
			&InventoryUpdated{Location: u(2), Product: u(3), OnHandChange: 3, OnHand: 3},
		},
		When: &GetLocInventoryReq{Location: u(1)},
		ThenResponse: &GetLocInventoryResp{
			Items: []*GetLocInventoryResp_Item{{Product: u(3), OnHand: 5, Available: 5}}},
	})

	Define(&Spec{
		Level: 2,
		Name:  "query inventory at root",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "Cola"},
			&ProductAdded{Uid: u(2), Sku: "Fanta"},
			&LocationAdded{Uid: u(3), Name: "Shelf", Parent: u(0)},
			&InventoryUpdated{Location: u(3), Product: u(2), OnHandChange: 2, OnHand: 2},
		},
		When: &GetLocInventoryReq{Location: u(0)},
		ThenResponse: &GetLocInventoryResp{
			Items: []*GetLocInventoryResp_Item{{Product: u(2), OnHand: 2, Available: 2}}},
	})

	container_with_gpus_inbound := []proto.Message{

		&ProductAdded{Uid: u(1), Sku: "NVidia 4080"},
		// we have a warehouse with unloading zone and a shelf
		&LocationAdded{Uid: u(1), Name: "Warehouse", Parent: u(0)},
		&LocationAdded{Uid: u(2), Name: "Unloading", Parent: u(1)},
		&LocationAdded{Uid: u(3), Name: "Shelf", Parent: u(1)},
		// 5 GPUS on a Shelf
		&InventoryUpdated{Location: u(3), Product: u(1), OnHandChange: 5, OnHand: 5},
		// we have a standalone container with some GPUs
		&LocationAdded{Uid: u(4), Name: "Container", Parent: u(0)},
		&InventoryUpdated{Location: u(4), Product: u(1), OnHandChange: 10, OnHand: 10},
		// container was moved to the unloading zone in warehouse
		&LocationMoved{Uid: u(4), NewParent: u(2)},
	}
	Define(&Spec{
		Level: 4,
		Name:  "moving container to warehouse increases total quantity",
		Given: container_with_gpus_inbound,
		// we query warehouse
		When: &GetLocInventoryReq{Location: u(1)},
		// warehouse should show 15 cards as being onHand
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: u(1), OnHand: 15, Available: 15},
		}},
	})

	Define(&Spec{
		Level: 4,
		Name:  "moving container to warehouse increases unloading quantity",
		Given: container_with_gpus_inbound,
		// we query unloading
		When: &GetLocInventoryReq{Location: u(2)},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: u(1), OnHand: 10, Available: 10},
		}},
	})

	Define(&Spec{
		Level: 4,
		Name:  "moving container with a reservation",
		Given: []proto.Message{
			// we have a warehouse with unloading zone and a shelf
			&LocationAdded{Uid: u(1), Name: "Warehouse", Parent: u(0)},
			&LocationAdded{Uid: u(2), Name: "Unloading", Parent: u(1)},
			&LocationAdded{Uid: u(3), Name: "Shelf", Parent: u(1)},
			&ProductAdded{Uid: u(4), Sku: "NVidia 4080"},
			// 5 GPUS on a Shelf
			&InventoryUpdated{Location: u(3), Product: u(4), OnHandChange: 5, OnHand: 5},
			// and 3 reserved
			&Reserved{
				Reservation: u(5),
				Code:        "sale1",
				Items:       []*Reserved_Item{{Product: u(4), Quantity: 3, Location: u(3)}},
			},
			// we have a standalone container with some GPUs
			&LocationAdded{Uid: u(6), Name: "Container", Parent: u(0)},
			&InventoryUpdated{Location: u(6), Product: u(4), OnHandChange: 10, OnHand: 10},
			// most of which was already promised to a customer
			&Reserved{
				Reservation: u(7),
				Code:        "sale3",
				Items:       []*Reserved_Item{{Product: u(4), Quantity: 9, Location: u(6)}},
			},
			// container was moved to the unloading zone in warehouse
			&LocationMoved{Uid: u(6), NewParent: u(2)},
		},
		When: &GetLocInventoryReq{Location: u(1)},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: u(4), OnHand: 15, Available: 3},
		}},
	})

	Define(&Spec{
		Level: 5,
		Name:  "reservation at location reduces availability at location",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "pixel"},
			&LocationAdded{Uid: u(2), Name: "Warehouse", Parent: u(0)},
			&InventoryUpdated{Location: u(2), Product: u(1), OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: u(3),
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 3, Location: u(2)}},
			},
		},
		When: &GetLocInventoryReq{Location: u(2)},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: u(1), OnHand: 10, Available: 7},
		}},
	})

	Define(&Spec{
		Level: 5,
		Name:  "cancelled reservation returns availability",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "pixel"},
			&LocationAdded{Uid: u(2), Name: "Warehouse", Parent: u(0)},
			&InventoryUpdated{Location: u(2), Product: u(1), OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: u(3),
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 3, Location: u(2)}},
			},
			&Cancelled{
				Reservation: u(3),
				Items:       []*Cancelled_Item{{Product: u(1), Released: 3, Location: u(2)}},
			},
		},
		When: &GetLocInventoryReq{Location: u(2)},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: u(1), OnHand: 10, Available: 10},
		}},
	})

	Define(&Spec{
		Level: 5,
		Name:  "reservation at location reduces availability globally",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "pixel"},
			&LocationAdded{Uid: u(2), Name: "Warehouse", Parent: u(0)},
			&InventoryUpdated{Location: u(2), Product: u(1), OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: u(3),
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 3, Location: u(2)}},
			},
		},
		When: &GetLocInventoryReq{Location: u(0)},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: u(1), OnHand: 10, Available: 7},
		}},
	})

	Define(&Spec{
		Level: 5,
		Name:  "multiple reservations stack",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "pixel"},
			&LocationAdded{Uid: u(2), Name: "Warehouse", Parent: u(0)},
			&InventoryUpdated{Location: u(2), Product: u(1), OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: u(3),
				Code:        "sale1",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 3, Location: u(2)}},
			},
			&Reserved{
				Reservation: u(4),
				Code:        "sale2",
				Items:       []*Reserved_Item{{Product: u(1), Quantity: 4, Location: u(2)}},
			},
		},
		When: &GetLocInventoryReq{Location: u(0)},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: u(1), OnHand: 10, Available: 3},
		}},
	})
}
