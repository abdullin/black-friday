package stock

import (
	. "black-friday/inventory/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

func init() {
	Define(&Spec{
		Name: "reserve sale with one item",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "GPU"},
			&LocationAdded{Id: 1, Name: "Shelf"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 10, OnHand: 10},
		},
		When: &ReserveReq{
			Reservation: "sale",
			Items: []*ReserveReq_Item{
				{Sku: "GPU", Quantity: 10},
			},
		},
		ThenResponse: &ReserveResp{
			Reservation: 1,
		},
		ThenEvents: []proto.Message{
			&Reserved{
				Reservation: 1,
				Code:        "sale",
				Items: []*Reserved_Item{
					{
						Product:  1,
						Quantity: 10,
					},
				},
			},
		},
	})

	Define(&Spec{
		Name: "reserve non-existent sku",
		When: &ReserveReq{
			Reservation: "test",
			Items: []*ReserveReq_Item{
				{Sku: "sale", Quantity: 1},
			},
		},
		ThenError: codes.NotFound,
	})

	Define(&Spec{
		Name: "experimental!",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "GPU"},
			&LocationAdded{Id: 1, Name: "WHS1"},
			&LocationAdded{Id: 2, Name: "WHS2"},
			&InventoryUpdated{Location: 1, Product: 1, OnHandChange: 5, OnHand: 5},
			&InventoryUpdated{Location: 2, Product: 1, OnHandChange: 5, OnHand: 5},
			&LambdaInstalled{
				Type: Lambda_RESERVE,
				Code: `
function reserve(order)
	channel = order:getTag("channel")
	if channel == 'store' then
		order:reserveFrom("WHS1")
	else
		order:reserveFrom("WHS2")
	end
end`,
			},
		},
		When: &ReserveReq{
			Reservation: "sale1",
			Items:       []*ReserveReq_Item{{Sku: "GPU", Quantity: 1}},
			Tags:        map[string]string{"channel": "online"},
		},
		ThenResponse: &ReserveResp{Reservation: 1},
		ThenEvents: []proto.Message{&Reserved{
			Reservation: 1,
			Code:        "sale1",
			Items:       []*Reserved_Item{{Product: 1, Quantity: 1, Location: 2}},
		}},
	})
}
