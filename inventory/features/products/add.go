package products

import (
	"black-friday/env/uid"
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"google.golang.org/grpc/status"
)

func Add(ctx fx.Tx, req *AddProductsReq) (r *AddProductsResp, status *status.Status) {

	id := ctx.GetSeq()

	results := make([]string, len(req.Skus))
	for i, sku := range req.Skus {

		id += 1
		uuid := uid.Str(id)
		e := &ProductAdded{Uid: uuid, Sku: sku}

		err, f := ctx.Apply(e)
		switch f {
		case fail.None:
		case fail.ConstraintUnique:
			return nil, ErrAlreadyExists
		default:
			return nil, ErrInternal(err, f)
		}

		results[i] = uuid
	}

	return &AddProductsResp{Uids: results}, nil
}
