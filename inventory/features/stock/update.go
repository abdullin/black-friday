package stock

import (
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"google.golang.org/grpc/status"
)

func Update(ctx fx.Tx, req *UpdateInventoryReq) (*UpdateInventoryResp, *status.Status) {

	var onHand int64

	ctx.QueryRow("SELECT OnHand FROM Inventory WHERE Location=? AND Product=?",
		req.Location, req.Product)(&onHand)

	onHand += req.OnHandChange

	if onHand < 0 {
		return nil, ErrNotEnough
	}

	e := &InventoryUpdated{
		Location:     req.Location,
		Product:      req.Product,
		OnHandChange: req.OnHandChange,
		OnHand:       onHand,
	}

	err, f := ctx.Apply(e)
	switch f {
	case fail.None:
	default:
		return nil, ErrInternal(err, f)
	}

	return &UpdateInventoryResp{OnHand: e.OnHand}, nil
}
