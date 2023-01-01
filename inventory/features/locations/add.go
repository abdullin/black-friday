package locations

import (
	"black-friday/env/uid"
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"google.golang.org/grpc/status"
)

func Add(c fx.Tx, req *AddLocationsReq) (*AddLocationsResp, *status.Status) {
	id := c.GetSeq("Locations")

	var addLoc func(parent string, ls []*AddLocationsReq_Loc) ([]*AddLocationsResp_Loc, *status.Status)

	addLoc = func(parent string, ls []*AddLocationsReq_Loc) ([]*AddLocationsResp_Loc, *status.Status) {

		var r []*AddLocationsResp_Loc

		for i, l := range ls {
			if l.Name == "" {
				return nil, ErrArgNil("name")
			}
			id += 1

			u := uid.Str(id)

			e := &LocationAdded{Name: l.Name, Uid: u, Parent: parent}
			node := &AddLocationsResp_Loc{Name: l.Name, Uid: u, Parent: parent}
			r = append(r, node)

			batch := i < len(ls)-1
			err, f := c.Apply(e, batch)
			switch f {
			case fail.None:
			case fail.ConstraintUnique:
				return nil, ErrAlreadyExists
			case fail.ConstraintForeign:
				return nil, ErrLocationNotFound
			default:
				return nil, ErrInternal(err, f)
			}

			children, st := addLoc(u, l.Locs)
			if err != nil {
				return nil, st
			}
			node.Locs = children
		}
		return r, nil
	}

	results, err := addLoc(req.Parent, req.Locs)
	if err != nil {
		return nil, err
	}

	return &AddLocationsResp{Locs: results}, nil
}
