package locations

import (
	"black-friday/fail"
	"black-friday/fx"
	. "black-friday/inventory/api"
)

func Add(c fx.Tx, req *AddLocationsReq) (*AddLocationsResp, error) {

	id := c.GetSeq("Locations")

	var addLoc func(parent int64, ls []*AddLocationsReq_Loc) ([]*AddLocationsResp_Loc, error)

	addLoc = func(parent int64, ls []*AddLocationsReq_Loc) ([]*AddLocationsResp_Loc, error) {

		var r []*AddLocationsResp_Loc

		for _, l := range ls {
			if l.Name == "" {
				return nil, ErrArgNil("name")
			}
			id += 1

			e := &LocationAdded{
				Name:   l.Name,
				Id:     id,
				Parent: parent,
			}
			node := &AddLocationsResp_Loc{
				Name:   l.Name,
				Id:     id,
				Parent: parent,
			}
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

			children, err := addLoc(id, l.Locs)
			if err != nil {
				return nil, err
			}
			node.Locs = children
		}
		return r, nil
	}

	results, err := addLoc(req.Parent, req.Locs)
	if err != nil {
		return nil, err
	}

	return &AddLocationsResp{
		Locs: results,
	}, nil

}
