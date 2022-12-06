package locations

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
		Level: 0,
		Name:  "add location",
		When: &AddLocationsReq{
			Locs: []*AddLocationsReq_Loc{
				{Name: "Shelf"},
			},
		},
		ThenResponse: &AddLocationsResp{Locs: []*AddLocationsResp_Loc{
			{Uid: u(1), Name: "Shelf"},
		}},
		ThenEvents: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Shelf"},
		},
	})

	Define(&Spec{
		Level: 1,
		Name:  "add locations to an existing one",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "WH"},
		},
		When: &AddLocationsReq{
			Parent: u(1),
			Locs: []*AddLocationsReq_Loc{
				{Name: "S1"},
				{Name: "S2"},
			},
		},
		ThenResponse: &AddLocationsResp{
			Locs: []*AddLocationsResp_Loc{
				{Uid: u(2), Name: "S1", Parent: u(1)},
				{Uid: u(3), Name: "S2", Parent: u(1)},
			},
		},
		ThenEvents: []proto.Message{
			&LocationAdded{Uid: u(2), Name: "S1", Parent: u(1)},
			&LocationAdded{Uid: u(3), Name: "S2", Parent: u(1)},
		},
	})
	Define(&Spec{
		Level: 1,
		Name:  "add nested locations to an existing one",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Warehouse"},
		},
		When: &AddLocationsReq{
			Parent: u(1),
			Locs: []*AddLocationsReq_Loc{
				{Name: "Shelf", Locs: []*AddLocationsReq_Loc{
					{Name: "Box"},
				}},
			},
		},
		ThenResponse: &AddLocationsResp{
			Locs: []*AddLocationsResp_Loc{
				{Uid: u(2), Name: "Shelf", Parent: u(1), Locs: []*AddLocationsResp_Loc{
					{Uid: u(3), Name: "Box", Parent: u(2)},
				}},
			},
		},
		ThenEvents: []proto.Message{
			&LocationAdded{Uid: u(2), Name: "Shelf", Parent: u(1)},
			&LocationAdded{Uid: u(3), Name: "Box", Parent: u(2)},
		},
	})

	Define(&Spec{
		Level: 1,
		Name:  "add location with wrong parent",
		Given: []proto.Message{},
		When: &AddLocationsReq{
			Parent: u(42),
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
			&LocationAdded{Uid: u(1), Name: "W"},
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
			&LocationAdded{Uid: u(1), Name: "WHS1"},
			&LocationAdded{Uid: u(2), Name: "Inbox", Parent: u(1)},
			&LocationAdded{Uid: u(3), Name: "WHS2"},
		},
		When: &AddLocationsReq{Parent: u(3), Locs: []*AddLocationsReq_Loc{{Name: "Inbox"}}},
		ThenResponse: &AddLocationsResp{Locs: []*AddLocationsResp_Loc{
			{Uid: u(4), Parent: u(3), Name: "Inbox"},
		}},
		ThenEvents: []proto.Message{&LocationAdded{Uid: u(4), Parent: u(3), Name: "Inbox"}},
	})

}
