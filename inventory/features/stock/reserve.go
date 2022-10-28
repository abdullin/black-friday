package stock

import (
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
)

func Reserve(a fx.Tx, r *ReserveReq) (*ReserveResp, error) {

	// by default, we reserve against the root.

	id := a.GetSeq("Reservations") + 1
	e := &Reserved{
		Reservation: id,
		Code:        r.Reservation,
	}

	for _, i := range r.Items {

		pid := a.LookupInt64("SELECT Id FROM Products WHERE Sku=?", i.Sku)
		if pid == 0 {
			return nil, ErrSkuNotFound(i.Sku)
		}
		e.Items = append(e.Items, &Reserved_Item{
			Product:  pid,
			Quantity: i.Quantity,
		})
	}

	err, f := a.Apply(e)
	switch f {
	case fail.None:
	default:
		return nil, ErrInternal(err, f)
	}

	return &ReserveResp{Reservation: id}, nil

}
