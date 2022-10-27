package products

import (
	"black-friday/fail"
	. "black-friday/inventory/api"
	"black-friday/inventory/app"
)

func Add(ctx *app.Context, req *AddProductsReq) (r *AddProductsResp, err error) {

	id := ctx.GetSeq("Products")

	results := make([]uint64, len(req.Skus))
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
			return nil, ErrDuplicateName
		default:
			return nil, ErrInternal(err, f)
		}

		results[i] = id
	}

	return &AddProductsResp{Ids: results}, nil
}
