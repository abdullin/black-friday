package locations

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {

	Define(&Spec{
		Level: 2,
		Name:  "query locations after removal",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Shelf"},
			&ProductAdded{Id: 1, Sku: "NVidia"},
			&InventoryUpdated{Product: 1, Location: 1, OnHandChange: 3, OnHand: 3},
			&InventoryUpdated{Product: 1, Location: 1, OnHandChange: -3, OnHand: 0},
		},
		When:         &GetLocInventoryReq{Location: 1},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{}},
	})

	Define(&Spec{
		Level: 2,
		Name:  "query one specific location",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Shelf1"},
		},
		When: &ListLocationsReq{Location: 1},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Id: 1, Name: "Shelf1"},
		}},
	})

	Define(&Spec{
		Level: 3,
		Name:  "query all locations in a tree",
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

	Define(&Spec{
		Level: 3,
		Name:  "query locations from another root",
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
	Define(&Spec{
		Level:     3,
		Name:      "query locations from non-existent location",
		When:      &ListLocationsReq{Location: 1},
		ThenError: ErrLocationNotFound,
	})

	Define(&Spec{
		Level: 2,
		Name:  "query all locations",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "WH1"},
			&LocationAdded{Id: 2, Name: "WH2"},
			&LocationAdded{Id: 3, Name: "Shelf", Parent: 1},
		},
		When: &ListLocationsReq{},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Name: "WH1", Id: 1, Chidren: []*ListLocationsResp_Loc{
				{Name: "Shelf", Id: 3, Parent: 1},
			}},
			{Name: "WH2", Id: 2},
		}},
	})

}
