package stock

import (
	"black-friday/fail"
	"black-friday/fx"
	"black-friday/inventory/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Update(ctx fx.Tx, req *api.UpdateInventoryReq) (r *api.UpdateInventoryResp, err error) {

	var onHand int64

	ctx.QueryRow("SELECT OnHand FROM Inventory WHERE Location=? AND Product=?",
		req.Location, req.Product)(&onHand)

	onHand += req.OnHandChange

	if onHand < 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "OnHand can't go negative!")
	}

	e := &api.InventoryUpdated{
		Location:     req.Location,
		Product:      req.Product,
		OnHandChange: req.OnHandChange,
		OnHand:       onHand,
	}

	err, f := ctx.Apply(e)
	switch f {
	case fail.None:
	default:
		return nil, api.ErrInternal(err, f)
	}

	return &api.UpdateInventoryResp{OnHand: e.OnHand}, nil
}
