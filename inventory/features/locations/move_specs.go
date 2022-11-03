package locations

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {

	Define(&Spec{
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

	Define(&Spec{
		Name: "recursive locations are not allowed",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Warehouse"},
			&LocationAdded{Id: 2, Name: "Container", Parent: 1},
		},
		When: &MoveLocationReq{
			Id:        1,
			NewParent: 2,
		},
		ThenError: ErrBadMove,
	})
	Define(&Spec{
		Name: "don't move location to itself",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Warehouse"},
		},
		When: &MoveLocationReq{
			Id:        1,
			NewParent: 1,
		},
		ThenError: ErrBadMove,
	})

	Define(&Spec{
		Name: "can't touch root",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Warehouse"},
		},
		When: &MoveLocationReq{
			Id:        0,
			NewParent: 1,
		},
		ThenError: ErrBadMove,
	})
}
