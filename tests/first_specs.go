package tests

import (
	"google.golang.org/grpc/codes"
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
			&LocationAdded{Id: 1, Name: "Shelf", Warehouse: 1},
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
		When: &ListLocationsReq{Warehouse: 1},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Location: 1, Name: "Shelf1", Warehouse: 1},
			{Location: 2, Name: "Shelf2", Warehouse: 1},
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

	register(&Spec{
		Name: "don't allow negative on-hand",
		Given: []proto.Message{
			&WarehouseCreated{Id: 1, Name: "WH1"},
			&LocationAdded{Id: 1, Name: "Shelf", Warehouse: 1},
			&ProductAdded{Id: 1, Sku: "NVidia"},
		},
		When:      &UpdateInventoryReq{Product: 1, Location: 1, OnHandChange: -1},
		ThenError: codes.FailedPrecondition,
	})

	register(&Spec{
		Name: "add locations",
		Given: []proto.Message{
			&WarehouseCreated{Id: 1, Name: "WH1"},
		},
		When: &AddLocationsReq{
			Warehouse: 1,
			Names:     []string{"L1", "L2"},
		},
		ThenResponse: &AddLocationsResp{
			Warehouse: 1,
			Ids:       []uint64{1, 2},
		},
		ThenEvents: []proto.Message{
			&LocationAdded{Warehouse: 1, Id: 1, Name: "L1"},
			&LocationAdded{Warehouse: 1, Id: 2, Name: "L2"},
		},
	})

	register(&Spec{
		Name:      "add locations without warehouse",
		Given:     []proto.Message{},
		When:      &AddLocationsReq{Warehouse: 42, Names: []string{"L1", "L2"}},
		ThenError: codes.FailedPrecondition,
	})

	register(&Spec{
		Name:      "add locations with zero warehouse id",
		Given:     []proto.Message{},
		When:      &AddLocationsReq{Warehouse: 0, Names: []string{"L1", "L2"}},
		ThenError: codes.InvalidArgument,
	})

	register(&Spec{
		Name: "query locations from another warehouse",
		Given: []proto.Message{
			&WarehouseCreated{Id: 1, Name: "WH1"},
			&WarehouseCreated{Id: 2, Name: "WH2"},
			&LocationAdded{Id: 1, Name: "Shelf", Warehouse: 1},
		},
		When:         &ListLocationsReq{Warehouse: 2},
		ThenResponse: &ListLocationsResp{},
	})
	register(&Spec{
		Name:         "query locations from non-existent warehouse",
		When:         &ListLocationsReq{Warehouse: 1},
		ThenResponse: &ListLocationsResp{},
	})

	/*
		register(&Spec{
			Name:      "insert duplicate warehouse name",
			When:      &CreateWarehouseReq{Names: []string{"W", "W"}},
			ThenError: codes.AlreadyExists,
		})

	*/

}
