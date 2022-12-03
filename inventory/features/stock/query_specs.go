package stock

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {
	Define(&Spec{
		Level: 3,
		Name:  "query inventory at a specific location",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "Cola"},
			&ProductAdded{Id: 2, Sku: "Fanta"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&InventoryUpdated{Location: 1, Product: 2, OnHandChange: 2, OnHand: 2},
		},
		When: &GetLocInventoryReq{Location: 1},
		ThenResponse: &GetLocInventoryResp{
			Items: []*GetLocInventoryResp_Item{{Product: 2, OnHand: 2, Available: 2}}},
	})

	Define(&Spec{
		Level: 3,
		Name:  "two boxes sum up their quantity at root",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "Epyc"},
			&LocationAdded{Id: 1, Name: "Shelf1"},
			&LocationAdded{Id: 2, Name: "Shelf2"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 2, OnHand: 2},
			&InventoryUpdated{Location: 2, Product: 1, OnHandChange: 3, OnHand: 3},
		},
		When: &GetLocInventoryReq{Location: 0},
		ThenResponse: &GetLocInventoryResp{
			Items: []*GetLocInventoryResp_Item{{Product: 1, OnHand: 5, Available: 5}}},
	})

	Define(&Spec{
		Level: 2,
		Name:  "query inventory at root",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "Cola"},
			&ProductAdded{Id: 2, Sku: "Fanta"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&InventoryUpdated{Location: 1, Product: 2, OnHandChange: 2, OnHand: 2},
		},
		When: &GetLocInventoryReq{Location: 0},
		ThenResponse: &GetLocInventoryResp{
			Items: []*GetLocInventoryResp_Item{{Product: 2, OnHand: 2, Available: 2}}},
	})

	container_with_gpus_inbound := []proto.Message{

		&ProductAdded{Id: 1, Sku: "NVidia 4080"},
		// we have a warehouse with unloading zone and a shelf
		&LocationAdded{Id: 1, Name: "Warehouse"},
		&LocationAdded{Id: 2, Name: "Unloading", Parent: 1},
		&LocationAdded{Id: 3, Name: "Shelf", Parent: 1},
		// 5 GPUS on a Shelf
		&InventoryUpdated{Location: 3, Product: 1, OnHandChange: 5, OnHand: 5},
		// we have a standalone container with some GPUs
		&LocationAdded{Id: 4, Name: "Container"},
		&InventoryUpdated{Location: 4, Product: 1, OnHandChange: 10, OnHand: 10},
		// container was moved to the unloading zone in warehouse
		&LocationMoved{Id: 4, NewParent: 2},
	}
	Define(&Spec{
		Level: 4,
		Name:  "moving container to warehouse increases total quantity",
		Given: container_with_gpus_inbound,
		// we query warehouse
		When: &GetLocInventoryReq{Location: 1},
		// warehouse should show 15 cards as being onHand
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: 1, OnHand: 15, Available: 15},
		}},
	})

	Define(&Spec{
		Level: 4,
		Name:  "moving container to warehouse increases unloading quantity",
		Given: container_with_gpus_inbound,
		// we query unloading
		When: &GetLocInventoryReq{Location: 2},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: 1, OnHand: 10, Available: 10},
		}},
	})

	Define(&Spec{
		Level: 4,
		Name:  "moving container with a reservation",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "NVidia 4080"},
			// we have a warehouse with unloading zone and a shelf
			&LocationAdded{Id: 1, Name: "Warehouse"},
			&LocationAdded{Id: 2, Name: "Unloading", Parent: 1},
			&LocationAdded{Id: 3, Name: "Shelf", Parent: 1},
			// 5 GPUS on a Shelf
			&InventoryUpdated{Location: 3, Product: 1, OnHandChange: 5, OnHand: 5},
			// and 3 reserved
			&Reserved{
				Reservation: 1,
				Code:        "sale1",
				Items:       []*Reserved_Item{{Product: 1, Quantity: 3, Location: 3}},
			},
			// we have a standalone container with some GPUs
			&LocationAdded{Id: 4, Name: "Container"},
			&InventoryUpdated{Location: 4, Product: 1, OnHandChange: 10, OnHand: 10},
			// most of which was already promised to a customer
			&Reserved{
				Reservation: 2,
				Code:        "sale3",
				Items:       []*Reserved_Item{{Product: 1, Quantity: 9, Location: 4}},
			},
			// container was moved to the unloading zone in warehouse
			&LocationMoved{Id: 4, NewParent: 2},
		},
		When: &GetLocInventoryReq{Location: 1},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: 1, OnHand: 15, Available: 3},
		}},
	})

	Define(&Spec{
		Level: 5,
		Name:  "reservation at location reduces availability at location",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "pixel"},
			&LocationAdded{Id: 1, Name: "Warehouse"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: 1,
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: 1, Quantity: 3, Location: 1}},
			},
		},
		When: &GetLocInventoryReq{Location: 1},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: 1, OnHand: 10, Available: 7},
		}},
	})

	Define(&Spec{
		Level: 5,
		Name:  "reservation at location reduces availability globally",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "pixel"},
			&LocationAdded{Id: 1, Name: "Warehouse"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: 1,
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: 1, Quantity: 3, Location: 1}},
			},
		},
		When: &GetLocInventoryReq{Location: 0},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: 1, OnHand: 10, Available: 7},
		}},
	})

	Define(&Spec{
		Level: 5,
		Name:  "multiple reservations stack",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "pixel"},
			&LocationAdded{Id: 1, Name: "Warehouse"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 10, OnHand: 10},
			&Reserved{
				Reservation: 1,
				Code:        "sale1",
				Items:       []*Reserved_Item{{Product: 1, Quantity: 3, Location: 1}},
			},
			&Reserved{
				Reservation: 2,
				Code:        "sale2",
				Items:       []*Reserved_Item{{Product: 1, Quantity: 4, Location: 1}},
			},
		},
		When: &GetLocInventoryReq{Location: 0},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: 1, OnHand: 10, Available: 3},
		}},
	})
}
