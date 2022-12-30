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

	// group all reservations by product
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

		stock := graphs.World.GetStock(int32(pid)).Clone()

		// remove all items
		for _, i := range is {

			qty, _ := stock.Update(int32(i.location), int32(-i.quantity), 0)

			e.Items = append(e.Items, &Fulfilled_Item{
				Product:  uid.Str(pid),
				Location: uid.Str(i.location),
				Removed:  i.quantity,
				OnHand:   int64(qty),
			})
		}
		// remove allocation

		for _, res := range reservation[pid] {
			_, reserved := stock.Update(int32(res.location), 0, int32(-res.quantity))

			e.Reserved = append(e.Reserved, &Fulfilled_Reserve{
				Product:  uid.Str(pid),
				Location: uid.Str(res.location),
				Quantity: res.quantity,
			})
			if reserved == 0 {
				return nil, status.Newf(codes.FailedPrecondition, "inventory not found")
			}
		}

		good := stock.IsValid()

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
