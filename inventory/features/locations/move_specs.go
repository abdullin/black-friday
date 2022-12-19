package locations

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {

	Define(&Spec{
		Level: 1,
		Name:  "move locations",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Warehouse", Parent: u(0)},
			&LocationAdded{Uid: u(2), Name: "Container", Parent: u(0)},
		},
		When:         &MoveLocationReq{Uid: u(2), NewParent: u(1)},
		ThenResponse: &MoveLocationResp{},
		ThenEvents: []proto.Message{
			&LocationMoved{Uid: u(2), OldParent: u(0), NewParent: u(1)},
		},
	})

	Define(&Spec{
		Level: 3,
		Name:  "recursive locations are not allowed",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Warehouse", Parent: u(0)},
			&LocationAdded{Uid: u(2), Name: "Container", Parent: u(1)},
		},
		When:      &MoveLocationReq{Uid: u(1), NewParent: u(2)},
		ThenError: ErrBadMove,
	})
	Define(&Spec{
		Level: 2,
		Name:  "don't move location to itself",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Warehouse", Parent: u(0)},
		},
		When:      &MoveLocationReq{Uid: u(1), NewParent: u(1)},
		ThenError: ErrBadMove,
	})

	Define(&Spec{
		Level: 1,
		Name:  "can't touch root",
		Given: []proto.Message{
			&LocationAdded{Uid: u(1), Name: "Warehouse", Parent: u(0)},
		},
		When:      &MoveLocationReq{Uid: u(0), NewParent: u(1)},
		ThenError: ErrBadMove,
	})
}
