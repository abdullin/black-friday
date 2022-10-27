package locations

import (
	"black-friday/fail"
	. "black-friday/inventory/api"
	"black-friday/inventory/app"
)

func Move(a *app.Context, r *MoveLocationReq) (*MoveLocationResp, error) {

	// root is not touchable
	if r.Id == 0 {
		return nil, ErrArgument
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
			return nil, ErrPrecondition
		}

		ancestor = a.LookupUint64("SELECT Parent FROM Locations WHERE Id=?", ancestor)

	}

	// this will be the old parent
	parent := a.LookupUint64("SELECT Parent FROM Locations WHERE Id=?", r.Id)

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
