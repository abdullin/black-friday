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

	locs, err := graphs.LoadLocTree(a)
	if err != nil {
		return nil, status.Convert(err)
	}

	skus := make(map[string]int64)

	for sku, quantity := range groups {
		var pid int64
		if !a.QueryRow("SELECT Id FROM Products WHERE Sku=?", sku)(&pid) {
			return nil, ErrProductNotFound
		}
		skus[sku] = pid

		inventory, err := graphs.LoadInventory(a, pid)
		if err != nil {
			return nil, status.Convert(err)
		}
		reserves, err := graphs.LoadReserves(a, pid)
		if err != nil {
			return nil, status.Convert(err)
		}

		reserves = append(reserves, graphs.Stock{
			Qty: quantity,
			Loc: loc,
		})

		enough := graphs.Resolves(a, locs, inventory, reserves)

		if !enough {
			return nil, status.Newf(codes.FailedPrecondition, "availability broken for product %d", pid)
		}
		/*
			// load products
			tree, err := graphs.LoadProductTree(a, pid)

			if err != nil {
				return nil, status.Convert(err)
			}

			_, _, found := graphs.Modify(tree, loc, 0, quantity)
			if !found {
				return nil, status.Newf(codes.FailedPrecondition, "no inventory for product %d", pid)
			}
			_, _, ok := graphs.Walk(tree)
			if !ok {
				return nil, status.Newf(codes.FailedPrecondition, "availability broken for product %d", pid)
			}

		*/

		e.Items = append(e.Items, &Reserved_Item{
			Product:  uid.Str(pid),
			Quantity: quantity,
			Location: r.Location,
		})
	}

	err, f := a.Apply(e)
	switch f {
	case fail.None:
	case fail.ConstraintUnique:
		return nil, ErrAlreadyExists
	default:
		return nil, ErrInternal(err, f)
	}

	return &ReserveResp{Reservation: uid.Str(id)}, nil

}
