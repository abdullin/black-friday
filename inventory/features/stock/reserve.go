package stock

import (
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"google.golang.org/grpc/status"
	"strings"
	"sync"
)

func Reserve(a fx.Tx, r *ReserveReq) (*ReserveResp, *status.Status) {

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

	for _, i := range r.Items {
		e.Items = append(e.Items, &Reserved_Item{
			Product:  skus[i.Sku],
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

var (
	cache *lua.FunctionProto
	lock  sync.Mutex
)

var usecache = true

func GetOrCompile(source string) (*lua.FunctionProto, error) {
	if usecache {

		lock.Lock()
		defer lock.Unlock()
		var err error
		if cache == nil {
			cache, err = CompileLua(source)
			if err != nil {
				return nil, err
			}
		}
		return cache, nil
	} else {
		return CompileLua(source)
	}

}

func DoLua(source string, L *lua.LState) error {
	p, err := GetOrCompile(source)
	if err != nil {
		return err
	}
	f := L.NewFunctionFromProto(p)
	L.Push(f)
	return L.PCall(0, lua.MultRet, nil)
}

func CompileLua(source string) (*lua.FunctionProto, error) {

	chunk, err := parse.Parse(strings.NewReader(source), "<string>")
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, "<string>")
	if err != nil {
		return nil, err
	}
	return proto, nil
}
