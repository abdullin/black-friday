package products

import (
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
)

func Add(ctx fx.Tx, req *AddProductsReq) (r *AddProductsResp, err error) {

	id := ctx.GetSeq("Products")

	results := make([]int64, len(req.Skus))
	for i, sku := range req.Skus {

		id += 1
		e := &ProductAdded{
			Id:  id,
			Sku: sku,
		}

		err, f := ctx.Apply(e)
		switch f {
		case fail.None:
		case fail.ConstraintUnique:
			return nil, ErrAlreadyExists
		default:
			return nil, ErrInternal(err, f)
		}

		results[i] = id
	}

	return &AddProductsResp{Ids: results}, nil
}
