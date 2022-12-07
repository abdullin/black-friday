package stock

import (
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"google.golang.org/grpc/status"
)

func Fulfill(ctx fx.Tx, req *FulfillReq) (*FulfillResp, *status.Status) {

	// temporary implementation
	var err error
	err, f := ctx.Apply(&Fulfilled{
		Reservation: req.Reservation,
		Items:       nil,
	})

	switch f {
	case fail.None:
	default:
		return nil, ErrInternal(err, f)
	}

	return &FulfillResp{}, nil

}
