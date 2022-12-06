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
			&LocationAdded{Uid: u(1), Name: "Shelf"},
			&ProductAdded{Uid: u(2), Sku: "NVidia"},
			&InventoryUpdated{Product: u(2), Location: u(1), OnHandChange: 3, OnHand: 3},
			&InventoryUpdated{Product: u(2), Location: u(1), OnHandChange: -3, OnHand: 0},
		},
		When:         &GetLocInventoryReq{Location: u(1)},
		ThenResponse: &GetLocInventoryResp{Items: []*GetLocInventoryResp_Item{}},
	})

	Define(&Spec{
		Level: 2,
		Name:  "query one specific location",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Shelf1"},
		},
		When: &ListLocationsReq{Location: u(1)},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Uid: u(1), Name: "Shelf1", Parent: u(0)},
		}},
	})

	Define(&Spec{
		Level: 3,
		Name:  "query all locations in a tree",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "WH"},
			&LocationAdded{Uid: u(2), Name: "Shelf1", Parent: u(1)},
			&LocationAdded{Uid: u(3), Name: "Shelf2", Parent: u(1)},
		},
		When: &ListLocationsReq{Location: u(1)},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Uid: u(1), Name: "WH", Parent: u(0), Chidren: []*ListLocationsResp_Loc{
				{Uid: u(2), Name: "Shelf1", Parent: u(1)},
				{Uid: u(3), Name: "Shelf2", Parent: u(1)},
			}},
		}},
	})

	Define(&Spec{
		Level: 3,
		Name:  "query locations from another root",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "WH1"},
			&LocationAdded{Uid: u(2), Name: "WH2"},
			&LocationAdded{Uid: u(3), Name: "Shelf", Parent: u(1)},
		},
		When: &ListLocationsReq{Location: u(2)},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{{
			Name:   "WH2",
			Uid:    u(2),
			Parent: u(0),
		}}},
	})
	Define(&Spec{
		Level:     3,
		Name:      "query locations from non-existent location",
		When:      &ListLocationsReq{Location: u(1)},
		ThenError: ErrLocationNotFound,
	})

	Define(&Spec{
		Level: 2,
		Name:  "query all locations",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "WH1"},
			&LocationAdded{Uid: u(2), Name: "WH2"},
			&LocationAdded{Uid: u(3), Name: "Shelf", Parent: u(1)},
		},
		When: &ListLocationsReq{},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Name: "WH1", Uid: u(1), Parent: u(0), Chidren: []*ListLocationsResp_Loc{
				{Name: "Shelf", Uid: u(3), Parent: u(1)},
			}},
			{Name: "WH2", Uid: u(2), Parent: u(0)},
		}},
	})

}
