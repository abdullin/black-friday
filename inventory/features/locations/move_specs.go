package locations

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {

	Define(&Spec{
		Level: 1,
		Comments: `
We can move locations around. For example, we could move a container
between the warehouses.
`,
		Name: "move locations",
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

	Define(&Spec{
		Name: "prevent moves that will break availability",
		Comments: `
The system should prevent moves that will break availability.

People will have to make a decision about what to do with that.

┌─ ── ── ── ── ── ── ─┐                     
 RESERVE: 2                                 
│ ┌─────────────────┐ │                     
│ │    Warehouse    │ │                     
  │ ┌─────────────┐ │   Move  ┌────────────┐
│ │ │Container: 2 ├─┼─┼──────▶│Container: 2│
│ │ └─────────────┘ │ │       └────────────┘
  └─────────────────┘                       
└─ ── ── ── ── ── ── ─┘                     
`,
		Level: 5,
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "GPU"},
			&LocationAdded{Uid: u(2), Name: "Warehouse", Parent: u(0)},
			&LocationAdded{Uid: u(3), Name: "Container", Parent: u(2)},
			&InventoryUpdated{Location: u(3), Product: u(1), OnHandChange: 2, OnHand: 2},
			&Reserved{
				Reservation: u(4),
				Code:        "sale",
				Items:       []*Reserved_Item{{Product: u(1), Location: u(2), Quantity: 2}},
			},
		},
		When: &MoveLocationReq{
			Uid:       u(3),
			NewParent: u(0),
		},
		ThenError: ErrBadMove,
	})

}
