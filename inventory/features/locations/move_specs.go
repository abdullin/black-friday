package locations

import (
	. "black-friday/inventory/api"
	"black-friday/specs"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func init() {

	specs.Add(&specs.S{
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

	specs.Add(&specs.S{
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
