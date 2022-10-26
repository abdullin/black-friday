package locations

import (
	"black-friday/fail"
	. "black-friday/inventory/api"
	"black-friday/inventory/app"
)

func Add(c *app.Context, req *AddLocationsReq) (*AddLocationsResp, error) {

	id := c.GetSeq("Locations")

	var addLoc func(parent uint64, ls []*AddLocationsReq_Loc) ([]*AddLocationsResp_Loc, error)

	addLoc = func(parent uint64, ls []*AddLocationsReq_Loc) ([]*AddLocationsResp_Loc, error) {

		var r []*AddLocationsResp_Loc

		for _, l := range ls {
			if l.Name == "" {
				return nil, ErrArgNil("Name")
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
				return nil, ErrDuplicateName
			case fail.ConstraintForeign:
				return nil, ErrNotFound
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