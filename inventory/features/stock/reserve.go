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

func Reserve(a fx.Tx, r *ReserveReq) (*ReserveResp, *status.Status) {

	// by default, we reserve against the root.

	if len(r.Items) == 0 {
		return nil, ErrArgument
	}

	id := a.GetSeq("Reservations") + 1
	e := &Reserved{
		Reservation: uid.Str(id),
		Code:        r.Reservation,
	}

	// group by sku
	groups := make(map[string]int64)
	for _, i := range r.Items {
		val, _ := groups[i.Sku]
		groups[i.Sku] = val + i.Quantity

	}

	loc := uid.Parse(r.Location)

	for sku, quantity := range groups {

		pid, found := graphs.World.SKUs[sku]

		if !found {
			return nil, ErrProductNotFound
		}

		stock := graphs.World.GetStock(pid).Clone()

		stock.Update(int32(loc), 0, int32(quantity))
		enough := stock.IsValid()

		if !enough {
			return nil, status.Newf(codes.FailedPrecondition, "availability broken for product %d", pid)
		}

		e.Items = append(e.Items, &Reserved_Item{
			Product:  uid.Str(int64(pid)),
			Quantity: quantity,
			Location: r.Location,
		})
	}

	err, f := a.Apply(e, false)
	switch f {
	case fail.None:
	case fail.ConstraintUnique:
		return nil, ErrAlreadyExists
	default:
		return nil, ErrInternal(err, f)
	}

	return &ReserveResp{Reservation: uid.Str(id)}, nil

}
