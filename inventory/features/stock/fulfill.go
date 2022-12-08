package stock

import (
	"black-friday/env/uid"
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"black-friday/inventory/features/graphs"
	"fmt"
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

	reservation := make(map[int64][]struct {
		location int64
		quantity int64
	})

	for rows.Next() {
		var product, location, quantity int64
		err := rows.Scan(&product, &location, &quantity)
		if err != nil {
			return nil, status.Convert(err)
		}
		res, _ := reservation[product]
		reservation[product] = append(res, struct {
			location int64
			quantity int64
		}{location: location, quantity: quantity})
	}

	if len(reservation) == 0 {
		return nil, status.Newf(codes.NotFound, "Reservation %d not found", rid)
	}

	// we need to ensure that we are:
	// fulfilling within our allocation (product is taken from location that is a child)
	// don't break any other allocations

	e := &Fulfilled{
		Reservation: req.Reservation,
	}

	// need to group products. product to location to onhand

	fill := make(map[int64][]struct {
		location, quantity int64
	})

	for _, i := range req.Items {
		pid := uid.Parse(i.Product)
		lid := uid.Parse(i.Location)
		s, _ := fill[pid]
		fill[pid] = append(s, struct{ location, quantity int64 }{location: lid, quantity: i.Quantity})
	}

	for pid, is := range fill {
		n, err := graphs.LoadProductTree(ctx, pid)
		if err != nil {
			return nil, status.Convert(err)
		}

		if fx.Explore {
			fmt.Printf("\nProduct tree at start: \n")
			graphs.Print(n, 2)
		}

		// remove all items
		for _, i := range is {

			on, _, found := graphs.Modify(n, i.location, -i.quantity, 0)
			if !found {
				return nil, status.Newf(codes.FailedPrecondition, "inventory not found")
			}
			e.Items = append(e.Items, &Fulfilled_Item{
				Product:  uid.Str(pid),
				Location: uid.Str(i.location),
				Removed:  i.quantity,
				OnHand:   on,
			})
		}
		// remove allocation

		for _, res := range reservation[pid] {
			_, _, found := graphs.Modify(n, res.location, 0, -res.quantity)
			if !found {
				return nil, status.Newf(codes.FailedPrecondition, "inventory not found")
			}
		}

		if fx.Explore {
			if fx.Explore {
				fmt.Printf("\nProduct tree at the end: \n")
				graphs.Print(n, 2)
			}

		}

		_, _, good := graphs.Walk(n)
		if !good {
			return nil, status.Newf(codes.FailedPrecondition, "broken availability")
		}

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
