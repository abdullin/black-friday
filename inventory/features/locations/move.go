package locations

import (
	"black-friday/env/uid"
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"google.golang.org/grpc/status"
)

func Move(a fx.Tx, r *MoveLocationReq) (*MoveLocationResp, *status.Status) {
	// root is not touchable

	id := uid.Parse(r.Uid)
	newParent := uid.Parse(r.NewParent)

	if id == 0 {
		return nil, ErrBadMove
	}
	// need to check if the new parent is not the child of the current node
	// OR the current node itself
	ancestor := newParent
	for {
		if ancestor == 0 {
			break
		}

		if ancestor == id {
			return nil, ErrBadMove
		}

		found := a.QueryRow("SELECT Parent FROM Locations WHERE Id=?", ancestor)(&ancestor)
		if !found {
			break
		}
	}
	// this will be the old parent
	var parent int64
	a.QueryRow("SELECT Parent FROM Locations WHERE Id=?", id)(&parent)

	err, f := a.Apply(&LocationMoved{
		Uid:       r.Uid,
		OldParent: uid.Str(parent),
		NewParent: r.NewParent,
	}, false)
	switch f {
	case fail.None:
	default:
		return nil, ErrInternal(err, f)

	}
	return &MoveLocationResp{}, nil
}
