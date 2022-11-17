package locations

import (
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
	"google.golang.org/grpc/status"
)

func Add(c fx.Tx, req *AddLocationsReq) (*AddLocationsResp, *status.Status) {
	id := c.GetSeq("Locations")

	var addLoc func(parent int64, ls []*AddLocationsReq_Loc) ([]*AddLocationsResp_Loc, *status.Status)

	addLoc = func(parent int64, ls []*AddLocationsReq_Loc) ([]*AddLocationsResp_Loc, *status.Status) {

		var r []*AddLocationsResp_Loc

		for _, l := range ls {
			if l.Name == "" {
				return nil, ErrArgNil("name")
			}
			id += 1

			e := &LocationAdded{Name: l.Name, Id: id, Parent: parent}
			node := &AddLocationsResp_Loc{Name: l.Name, Id: id, Parent: parent}
			r = append(r, node)

			err, f := c.Apply(e)
			switch f {
			case fail.None:
			case fail.ConstraintUnique:
				return nil, ErrAlreadyExists
			case fail.ConstraintForeign:
				return nil, ErrLocationNotFound
			default:
				return nil, ErrInternal(err, f)
			}

			children, st := addLoc(id, l.Locs)
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
