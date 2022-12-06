package products

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
		Level:        0,
		Name:         "add new products",
		When:         &AddProductsReq{Skus: []string{"one"}},
		ThenResponse: &AddProductsResp{Uids: []string{u(1)}},
		ThenEvents: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "one"},
		},
	})

	Define(&Spec{
		Level:        0,
		Name:         "add new products",
		When:         &AddProductsReq{Skus: []string{"one", "two"}},
		ThenResponse: &AddProductsResp{Uids: []string{u(1), u(2)}},
		ThenEvents: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "one"},
			&ProductAdded{Uid: u(2), Sku: "two"},
		},
	})

	Define(&Spec{
		Level: 1,
		Name:  "can't add product with duplicate skus",
		Given: []proto.Message{
			&ProductAdded{Uid: u(1), Sku: "cola"},
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
