package stock

import (
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	lua "github.com/yuin/gopher-lua"
)

func Reserve(a fx.Tx, r *ReserveReq) (*ReserveResp, error) {

	// by default, we reserve against the root.

	id := a.GetSeq("Reservations") + 1
	e := &Reserved{
		Reservation: id,
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

	var code string
	found := a.QueryRow("SELECT Code FROM Lambdas WHERE Type=?", Lambda_RESERVE.String())(&code)
	if found {
		vm := lua.NewState()
		defer vm.Close()
		// we have a custom handler
		// let's run our reservation against it

		tags := vm.NewTable()
		for k, v := range r.Tags {
			tags.RawSetString(k, lua.LString(v))
		}

		order := vm.NewTable()
		order.RawSetString("tags", tags)
		order.RawSetString("id", lua.LString(r.Reservation))

		items := vm.NewTable()

		for _, v := range r.Items {
			it := vm.NewTable()
			it.RawSetString("id", lua.LNumber(skus[v.Sku]))
			it.RawSetString("sku", lua.LString(v.Sku))
			it.RawSetString("quantity", lua.LNumber(v.Quantity))
			items.Append(it)
		}

		order.RawSetString("items", items)

		vm.SetGlobal("order", order)

		ReserveAll := func(s *lua.LState) int {
			location := s.ToString(1) /* get argument */
			var locId int64
			found := a.QueryRow("SELECT Id FROM Locations WHERE Name=?", location)(&locId)
			if !found {
				s.RaiseError("location not found")
				return 0
			}

			items.ForEach(func(key lua.LValue, v lua.LValue) {
				t := v.(*lua.LTable)

				id := int64(t.RawGetString("id").(lua.LNumber))
				quantity := int64(t.RawGetString("quantity").(lua.LNumber))
				e.Items = append(e.Items, &Reserved_Item{
					Product:  id,
					Quantity: quantity,
					Location: locId,
				})
			})
			return 0 /* number of results */
		}

		vm.SetGlobal("reserveAll", vm.NewFunction(ReserveAll))

		err := vm.DoString(code)
		if err != nil {
			return nil, err
		}

	} else {
		for _, i := range r.Items {
			e.Items = append(e.Items, &Reserved_Item{
				Product:  skus[i.Sku],
				Quantity: i.Quantity,
			})
		}
	}

	err, f := a.Apply(e)
	switch f {
	case fail.None:
	default:
		return nil, ErrInternal(err, f)
	}

	return &ReserveResp{Reservation: id}, nil

}
