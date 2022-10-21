package tests

import (
	. "black-friday/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func init() {

	register(&Spec{
		Name: "query inventory",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "Cola"},
			&ProductAdded{Id: 2, Sku: "Fanta"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&InventoryUpdated{Location: 1, Product: 2, OnHandChange: 2, OnHand: 2},
		},
		When:         &GetLocInventoryReq{Location: 1},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{{Product: 2, OnHand: 2}}},
	})

	register(&Spec{
		Name: "query one specific location",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Shelf1"},
		},
		When: &ListLocationsReq{Location: 1},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Id: 1, Name: "Shelf1"},
		}},
	})

	register(&Spec{
		Name: "query all locations in a tree",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "WH"},
			&LocationAdded{Id: 2, Name: "Shelf1", Parent: 1},
			&LocationAdded{Id: 3, Name: "Shelf2", Parent: 1},
		},
		When: &ListLocationsReq{Location: 1},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Id: 1, Name: "WH", Chidren: []*ListLocationsResp_Loc{
				{Id: 2, Name: "Shelf1", Parent: 1},
				{Id: 3, Name: "Shelf2", Parent: 1},
			}},
		}},
	})

	register(&Spec{
		Name: "query locations after removal",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Shelf"},
			&ProductAdded{Id: 1, Sku: "NVidia"},
			&InventoryUpdated{Product: 1, Location: 1, OnHandChange: 3, OnHand: 3},
			&InventoryUpdated{Product: 1, Location: 1, OnHandChange: -3, OnHand: 0},
		},
		When:         &GetLocInventoryReq{Location: 1},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{}},
	})

	register(&Spec{
		Name: "don't allow negative on-hand",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Shelf"},
			&ProductAdded{Id: 1, Sku: "NVidia"},
		},
		When:      &UpdateInventoryReq{Product: 1, Location: 1, OnHandChange: -1},
		ThenError: codes.FailedPrecondition,
	})

	register(&Spec{
		Name: "add locations to an existing one",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "WH"},
		},
		When: &AddLocationsReq{
			Parent: 1,
			Locs: []*AddLocationsReq_Loc{
				{Name: "S1"},
				{Name: "S2"},
			},
		},
		ThenResponse: &AddLocationsResp{
			Locs: []*AddLocationsResp_Loc{
				{Id: 2, Name: "S1", Parent: 1},
				{Id: 3, Name: "S2", Parent: 1},
			},
		},
		ThenEvents: []proto.Message{
			&LocationAdded{Id: 2, Name: "S1", Parent: 1},
			&LocationAdded{Id: 3, Name: "S2", Parent: 1},
		},
	})

	register(&Spec{
		Name: "add nested locations to an existing one",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Warehouse"},
		},
		When: &AddLocationsReq{
			Parent: 1,
			Locs: []*AddLocationsReq_Loc{
				{Name: "Shelf", Locs: []*AddLocationsReq_Loc{
					{Name: "Box"},
				}},
			},
		},
		ThenResponse: &AddLocationsResp{
			Locs: []*AddLocationsResp_Loc{
				{Id: 2, Name: "Shelf", Parent: 1, Locs: []*AddLocationsResp_Loc{
					{Id: 3, Name: "Box", Parent: 2},
				}},
			},
		},
		ThenEvents: []proto.Message{
			&LocationAdded{Id: 2, Name: "Shelf", Parent: 1},
			&LocationAdded{Id: 3, Name: "Box", Parent: 2},
		},
	})

	register(&Spec{
		Name:  "add location with wrong parent",
		Given: []proto.Message{},
		When: &AddLocationsReq{
			Parent: 42,
			Locs: []*AddLocationsReq_Loc{{
				Name: "L",
			}},
		},
		ThenError: codes.NotFound,
	})

	register(&Spec{
		Name:      "add location with nill name",
		Given:     []proto.Message{},
		When:      &AddLocationsReq{Locs: []*AddLocationsReq_Loc{{}}},
		ThenError: codes.InvalidArgument,
	})

	register(&Spec{
		Name: "query locations from another root",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "WH1"},
			&LocationAdded{Id: 2, Name: "WH2"},
			&LocationAdded{Id: 3, Name: "Shelf", Parent: 1},
		},
		When: &ListLocationsReq{Location: 2},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{{
			Name: "WH2",
			Id:   2,
		}}},
	})
	register(&Spec{
		Name:      "query locations from non-existent location",
		When:      &ListLocationsReq{Location: 1},
		ThenError: codes.NotFound,
	})

	register(&Spec{
		Name: "insert duplicate location name in a batch",
		When: &AddLocationsReq{Locs: []*AddLocationsReq_Loc{
			{Name: "W"},
			{Name: "W"},
		}},
		ThenError: codes.AlreadyExists,
	})

	register(&Spec{
		Name: "add location with duplicate name",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "W"},
		},
		When: &AddLocationsReq{Locs: []*AddLocationsReq_Loc{
			{Name: "W"},
		}},
		ThenError: codes.AlreadyExists,
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
	register(&Spec{
		Name:  "moving container to warehouse increases total quantity",
		Given: container_with_gpus_inbound,
		// we query warehouse
		When: &GetLocInventoryReq{Location: 1},
		// warehouse should show 15 cards as being onHand
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: 1, OnHand: 15},
		}},
	})

	register(&Spec{
		Name:  "moving container to warehouse increases unloading quantity",
		Given: container_with_gpus_inbound,
		// we query unloading
		When: &GetLocInventoryReq{Location: 2},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{
			{Product: 1, OnHand: 10},
		}},
	})

	register(&Spec{
		Name: "move locations",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Warehouse"},
			&LocationAdded{Id: 2, Name: "Container"},
		},
		When: &MoveLocationReq{
			Id:        2,
			NewParent: 1,
		},
		ThenResponse: &MoveLocationResp{},
		ThenEvents: []proto.Message{
			&LocationMoved{Id: 2, OldParent: 0, NewParent: 1},
		},
	})

	register(&Spec{
		Name: "recursive locations are not allowed",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Warehouse"},
			&LocationAdded{Id: 2, Name: "Container", Parent: 1},
		},
		When: &MoveLocationReq{
			Id:        1,
			NewParent: 2,
		},
		ThenError: codes.FailedPrecondition,
	})

}
