package tests

import (
	"google.golang.org/protobuf/proto"
	. "sdk-go/protos"
)

func init() {

	register(&Spec{
		Name: "query inventory",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "Cola"},
			&ProductAdded{Id: 2, Sku: "Fanta"},
			&WarehouseCreated{Id: 1, Name: "WH1"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&InventoryUpdated{Location: 1, Product: 2, OnHandChange: 2, OnHand: 2},
		},
		When:         &GetInventoryReq{Location: 1},
		ThenResponse: &GetInventoryResp{Items: []*GetInventoryResp_Item{{Product: 2, OnHand: 2}}},
	})

	register(&Spec{
		Name: "query locations",
		Given: []proto.Message{

			&WarehouseCreated{Id: 1, Name: "WH1"},
			&LocationAdded{Id: 1, Name: "Shelf1", Warehouse: 1},
			&LocationAdded{Id: 2, Name: "Shelf2", Warehouse: 1},
		},
		When: &ListLocationsReq{},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Id: 1, Name: "Shelf1"},
			{Id: 2, Name: "Shelf2"},
		}},
	})

	register(&Spec{
		Name: "query locations after removal",
		Given: []proto.Message{

			&WarehouseCreated{Id: 1, Name: "WH1"},
			&LocationAdded{Id: 1, Name: "Shelf", Warehouse: 1},
			&ProductAdded{Id: 1, Sku: "NVidia"},
			&InventoryUpdated{Product: 1, Location: 1, OnHandChange: 3, OnHand: 3},
			&InventoryUpdated{Product: 1, Location: 1, OnHandChange: -3, OnHand: 0},
		},
		When:         &GetInventoryReq{Location: 1},
		ThenResponse: &GetInventoryResp{Items: []*GetInventoryResp_Item{}},
	})

}
