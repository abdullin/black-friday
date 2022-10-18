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
				{Id: 1, Name: "Shelf1"},
				{Id: 2, Name: "Shelf2"},
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
			&LocationAdded{Id: 1, Name: "WH"},
		},
		When: &AddLocationsReq{
			Parent: 1,
			Locs: []*AddLocationsReq_Loc{
				{Name: "S1", Locs: []*AddLocationsReq_Loc{
					{Name: "S2"},
				}},
			},
		},
		ThenResponse: &AddLocationsResp{
			Locs: []*AddLocationsResp_Loc{
				{Id: 2, Name: "S1", Parent: 1, Locs: []*AddLocationsResp_Loc{
					{Id: 3, Name: "S2", Parent: 2},
				}},
			},
		},
		ThenEvents: []proto.Message{
			&LocationAdded{Id: 2, Name: "S1", Parent: 1},
			&LocationAdded{Id: 3, Name: "S2", Parent: 2},
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

}
