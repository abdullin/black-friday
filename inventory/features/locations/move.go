package locations

import (
	. "black-friday/inventory/api"
	"black-friday/inventory/app"
	"database/sql"
)

func Move(a *app.Context, r *MoveLocationReq) (*MoveLocationResp, error) {
	parent, err := a.QueryUint64("SELECT Parent FROM Locations WHERE Id=?", r.Id)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	a.Apply(&LocationMoved{
		Id:        r.Id,
		OldParent: parent,
		NewParent: r.NewParent,
	})

	return &MoveLocationResp{}, nil

}