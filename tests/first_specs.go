package tests

import (
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
			&QuantityUpdated{Location: 1, Product: 2, Quantity: 1, Total: 2, Before: 0},
		},
		When:         &GetInventoryReq{Location: 1},
		ThenResponse: &GetInventoryResp{Items: []*GetInventoryResp_Item{{Product: 2, Quantity: 2}}},
	})

	register(&Spec{
		Name: "query locations",
		Given: []proto.Message{
			&LocationAdded{Id: 1, Name: "Shelf1"},
			&LocationAdded{Id: 2, Name: "Shelf2"},
		},
		When: &ListLocationsReq{},
		ThenResponse: &ListLocationsResp{Locs: []*ListLocationsResp_Loc{
			{Id: 1, Name: "Shelf1"},
			{Id: 2, Name: "Shelf2"},
		}},
	})

}
