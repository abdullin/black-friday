package products

import (
	. "black-friday/inventory/api"
	"google.golang.org/protobuf/proto"
)

func init() {
	Define(&Spec{
		Name:         "create new products",
		When:         &AddProductsReq{Skus: []string{"one", "two"}},
		ThenResponse: &AddProductsResp{Ids: []int64{1, 2}},
		ThenEvents: []proto.Message{
			&ProductAdded{Id: 1, Sku: "one"},
			&ProductAdded{Id: 2, Sku: "two"},
		},
	})

	Define(&Spec{
		Name:      "one failing product fails the batch",
		When:      &AddProductsReq{Skus: []string{"one", "one"}},
		ThenError: ErrAlreadyExists,
	})
}
