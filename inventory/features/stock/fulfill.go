package stock

import (
	"black-friday/env/uid"
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"google.golang.org/grpc/status"
)

func Fulfill(ctx fx.Tx, req *FulfillReq) (*FulfillResp, *status.Status) {

	// load reservation details

	rid := uid.Parse(req.Reservation)

	rows, err := ctx.QueryHack(`SELECT Product, Location, Quantity FROM Reserves WHERE Reservation=?`, rid)
	if err != nil {
		return nil, status.Convert(err)
	}

	defer rows.Close()

	lookup := make(map[int64]struct {
		location int64;
		quantity int64
	})

	for rows.Next() {
		var product, location, quantity int64
		err := rows.Scan(&product, &location, &quantity)
		if err != nil {
			return nil, status.Convert(err)
		}
		lookup[product] = struct {
			location int64
			quantity int64
		}{location: location, quantity: quantity}
	}

	if len(lookup) == 0 {
		return nil, ErrReservationNotFound
	}

	// temporary implementation
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
