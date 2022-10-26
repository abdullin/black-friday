package products

import (
	"black-friday/fail"
	"black-friday/inventory/api"
	"black-friday/inventory/app"
)

func Add(ctx *app.Context, req *api.AddProductsReq) (r *api.AddProductsResp, err error) {

	id := ctx.GetSeq("Products")

	results := make([]uint64, len(req.Skus))
	for i, sku := range req.Skus {

		id += 1
		e := &api.ProductAdded{
			Id:  id,
			Sku: sku,
		}

		err, f := ctx.Apply(e)
		switch f {
		case fail.None:
		case fail.ConstraintUnique:
			return nil, api.ErrDuplicateName
		default:
			return nil, api.ErrInternal(err, f)
		}

		results[i] = id
	}

	return &api.AddProductsResp{Ids: results}, nil
}
