package stock

import (
	"black-friday/env/uid"
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"black-friday/inventory/features/graphs"
	"google.golang.org/grpc/codes"
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
		location int64
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

	// we need to ensure that we are:
	// fulfilling within our allocation (product is taken from location that is a child)
	// don't break any other allocations

	e := &Fulfilled{
		Reservation: req.Reservation,
	}

	for _, i := range req.Items {
		n, err := graphs.LoadProductTree(ctx, uid.Parse(i.Product))
		if err != nil {
			return nil, status.Convert(err)
		}

		on, _, found := graphs.Modify(n, uid.Parse(i.Location), -i.Quantity, -i.Quantity)
		if !found {
			return nil, status.Newf(codes.FailedPrecondition, "inventory not found")
		}
		_, _, good := graphs.Walk(n)
		if !good {
			return nil, status.Newf(codes.FailedPrecondition, "broken availability")
		}
		e.Items = append(e.Items, &Fulfilled_Item{
			Product:  i.Product,
			Location: i.Location,
			Removed:  i.Quantity,
			OnHand:   on,
		})

	}

	// temporary implementation
	err, f := ctx.Apply(e)

	switch f {
	case fail.None:
	default:
		return nil, ErrInternal(err, f)
	}

	return &FulfillResp{}, nil

}
