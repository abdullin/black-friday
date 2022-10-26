package stock

import (
	. "black-friday/inventory/api"
	"black-friday/specs"
	"google.golang.org/protobuf/proto"
)

func init() {
	specs.Add(&specs.S{
		Name: "query inventory",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "Cola"},
			&ProductAdded{Id: 2, Sku: "Fanta"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&InventoryUpdated{Location: 1, Product: 2, OnHandChange: 2, OnHand: 2},
		},
		When: &GetLocInventoryReq{Location: 1},
		ThenResponse: &GetLocInventoryResp{
			Items: []*GetLocInventoryResp_Item{{Product: 2, OnHand: 2}}},
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
	specs.Add(&specs.S{
		Name:  "moving container to warehouse increases total quantity",
		Given: container_with_gpus_inbound,
		// we query warehouse
		When: &GetLocInventoryReq{Location: 1},
		// warehouse should show 15 cards as being onHand
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: 1, OnHand: 15},
		}},
	})

	specs.Add(&specs.S{
		Name:  "moving container to warehouse increases unloading quantity",
		Given: container_with_gpus_inbound,
		// we query unloading
		When: &GetLocInventoryReq{Location: 2},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: 1, OnHand: 10},
		}},
	})
}
