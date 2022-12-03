package locations

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {

	Define(&Spec{
		Level: 0,
		Name:  "add location",
		When: &AddLocationsReq{
			Locs: []*AddLocationsReq_Loc{
				{Name: "Shelf"},
			},
		},
		ThenResponse: &AddLocationsResp{Locs: []*AddLocationsResp_Loc{
			{Id: 1, Name: "Shelf"},
		}},
		ThenEvents: []proto.Message{
			&LocationAdded{Id: 1, Name: "Shelf"},
		},
	})

	Define(&Spec{
		Level: 1,
		Name:  "add locations to an existing one",
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
	Define(&Spec{
		Level: 1,
		Name:  "add nested locations to an existing one",
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

	Define(&Spec{
		Level: 1,
		Name:  "add location with wrong parent",
		Given: []proto.Message{},
		When: &AddLocationsReq{
			Parent: 42,
			Locs: []*AddLocationsReq_Loc{{
				Name: "L",
			}},
		},
		ThenError: ErrLocationNotFound,
	})

	Define(&Spec{
		Level:     0,
		Name:      "add location with nil name",
		Given:     []proto.Message{},
		When:      &AddLocationsReq{Locs: []*AddLocationsReq_Loc{{}}},
		ThenError: ErrArgNil("name"),
	})
	Define(&Spec{
		Level: 2,
		Name:  "insert duplicate location name in a batch",
		When: &AddLocationsReq{Locs: []*AddLocationsReq_Loc{
			{Name: "W"},
			{Name: "W"},
		}},
		ThenError: ErrAlreadyExists,
	})

	Define(&Spec{
		Level: 2,
		Name:  "add location with duplicate name",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "W"},
		},
		When: &AddLocationsReq{Locs: []*AddLocationsReq_Loc{
			{Name: "W"},
		}},
		ThenError: ErrAlreadyExists,
	})
	Define(&Spec{
		Level: 2,
		Name:  "duplicates are OK, if they don't share a parent",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "WHS1"},
			&LocationAdded{Id: 2, Name: "Inbox", Parent: 1},
			&LocationAdded{Id: 3, Name: "WHS2"},
		},
		When: &AddLocationsReq{Parent: 3, Locs: []*AddLocationsReq_Loc{{Name: "Inbox"}}},
		ThenResponse: &AddLocationsResp{Locs: []*AddLocationsResp_Loc{
			{Id: 4, Parent: 3, Name: "Inbox"},
		}},
		ThenEvents: []proto.Message{&LocationAdded{Id: 4, Parent: 3, Name: "Inbox"}},
	})

}
