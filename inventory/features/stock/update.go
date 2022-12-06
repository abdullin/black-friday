package stock

import (
	"black-friday/env/uid"
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"google.golang.org/grpc/status"
)

func Update(ctx fx.Tx, req *UpdateInventoryReq) (*UpdateInventoryResp, *status.Status) {

	var onHand int64

	lid := uid.Parse(req.Location)
	pid := uid.Parse(req.Product)

	ctx.QueryRow("SELECT OnHand FROM Inventory WHERE Location=? AND Product=?", lid, pid)(&onHand)

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
