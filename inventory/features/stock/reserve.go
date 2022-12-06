package stock

import (
	"black-friday/env/uid"
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Reserve(a fx.Tx, r *ReserveReq) (*ReserveResp, *status.Status) {

	// by default, we reserve against the root.

	id := a.GetSeq("Reservations") + 1
	e := &Reserved{
		Reservation: uid.Str(id),
		Code:        r.Reservation,
	}

	skus := make(map[string]int64)

	for _, r := range r.Items {
		var pid int64
		if !a.QueryRow("SELECT Id FROM Products WHERE Sku=?", r.Sku)(&pid) {
			return nil, ErrProductNotFound
		}
		skus[r.Sku] = pid
	}

	// this is a slow route for now

	res, st := Query(a, &GetLocInventoryReq{Location: r.Location})
	if st.Code() != codes.OK {
		return nil, st
	}

	available := make(map[int64]int64)

	for _, prod := range res.Items {
		available[uid.Parse(prod.Product)] = prod.Available
	}

	for _, i := range r.Items {
		productId := skus[i.Sku]

		available, _ := available[productId]
		if available < i.Quantity {
			return nil, ErrNotEnough
		}

		e.Items = append(e.Items, &Reserved_Item{
			Product:  uid.Str(productId),
			Quantity: i.Quantity,
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

	// here is the tricky part. We need to walk the hierarchy to see if things are still good

	scope := uid.Parse(r.Location)
	for scope != 0 {

		var parent int64

		a.QueryRow("SELECT Parent FROM Locations WHERE Id=?", scope)(&parent)
		res, st := Query(a, &GetLocInventoryReq{Location: uid.Str(scope)})
		if st.Code() != codes.OK {
			return nil, st
		}
		// checking availability

		for _, i := range res.Items {
			if i.Available < 0 {
				// we broke some constraint!
				return nil, ErrNotEnough
			}
		}
		scope = parent
	}

	return &ReserveResp{Reservation: uid.Str(id)}, nil

}
