package locations

import (
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
)

func Move(a fx.Tx, r *MoveLocationReq) (*MoveLocationResp, error) {

	// root is not touchable
	if r.Id == 0 {
		return nil, ErrBadMove
	}

	// need to check if the new parent is not the child of the current node
	// OR the current node itself
	// we shouldn't let the new parent to be a child of the node being moved
	// so let's inspect that

	ancestor := r.NewParent
	for {
		if ancestor == 0 {
			break
		}

		if ancestor == r.Id {
			return nil, ErrBadMove
		}

		found := a.QueryRow("SELECT Parent FROM Locations WHERE Id=?", ancestor)(&ancestor)
		if !found {
			break
		}

	}

	// this will be the old parent
	var parent int64
	a.QueryRow("SELECT Parent FROM Locations WHERE Id=?", r.Id)(&parent)

	err, f := a.Apply(&LocationMoved{
		Id:        r.Id,
		OldParent: parent,
		NewParent: r.NewParent,
	})
	switch f {
	case fail.None:
	default:
		return nil, ErrInternal(err, f)

	}

	return &MoveLocationResp{}, nil

}
