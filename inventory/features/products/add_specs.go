package products

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {
	Define(&Spec{
		Level:        0,
		Name:         "add new products",
		When:         &AddProductsReq{Skus: []string{"one"}},
		ThenResponse: &AddProductsResp{Ids: []int64{1}},
		ThenEvents: []proto.Message{
			&ProductAdded{Id: 1, Sku: "one"},
		},
	})

	Define(&Spec{
		Level:        0,
		Name:         "add new products",
		When:         &AddProductsReq{Skus: []string{"one", "two"}},
		ThenResponse: &AddProductsResp{Ids: []int64{1, 2}},
		ThenEvents: []proto.Message{
			&ProductAdded{Id: 1, Sku: "one"},
			&ProductAdded{Id: 2, Sku: "two"},
		},
	})

	Define(&Spec{
		Level: 1,
		Name:  "can't add product with duplicate skus",
		Given: []proto.Message{
			&ProductAdded{Id: 1, Sku: "cola"},
		},
		When:      &AddProductsReq{Skus: []string{"cola"}},
		ThenError: ErrAlreadyExists,
	})

	Define(&Spec{
		Level:     1,
		Name:      "can't add multiple product with duplicate skus",
		When:      &AddProductsReq{Skus: []string{"one", "one"}},
		ThenError: ErrAlreadyExists,
	})
}
