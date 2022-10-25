package inventory

import (
	. "black-friday/api"
	"black-friday/fail"
	"context"
	"database/sql"
)

func (s *App) MoveLocation(ctx context.Context, req *MoveLocationReq) (*MoveLocationResp, error) {
	return nil, ErrNotUnimplemented
}

func (s *App) ListLocations(ctx context.Context, req *ListLocationsReq) (*ListLocationsResp, error) {

	tx := s.GetTx(ctx)

	rows, err := tx.Tx.QueryContext(ctx, `
WITH RECURSIVE cte_Locations(Id, Parent, Name) AS (
	SELECT l.Id, l.Parent, l.Name
	FROM Locations l
	WHERE l.Id = ?

	UNION ALL

	SELECT l.Id, l.Parent, l.Name
	FROM Locations l
	JOIN cte_Locations c ON c.Id = l.Parent
)
SELECT * FROM cte_Locations
`, zeroToNill(req.Location))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*ListLocationsResp_Loc

	lookup := make(map[uint64]*ListLocationsResp_Loc)

	for rows.Next() {
		var id uint64
		var parent sql.NullInt64
		var name string
		err := rows.Scan(&id, &parent, &name)
		if err != nil {
			return nil, err
		}

		loc := &ListLocationsResp_Loc{
			Name:    name,
			Id:      id,
			Parent:  uint64(parent.Int64),
			Chidren: nil,
		}
		lookup[id] = loc
		if parent, found := lookup[loc.Parent]; found {
			parent.Chidren = append(parent.Chidren, loc)
		} else {
			results = append(results, loc)
		}
	}
	// we should at list get one location
	if len(results) == 0 {
		return nil, ErrNotFound
	}
	return &ListLocationsResp{Locs: results}, nil
}

func (s *App) AddLocations(ctx context.Context, req *AddLocationsReq) (r *AddLocationsResp, e error) {

	tx := s.GetTx(ctx)

	id := tx.GetSeq("Locations")

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

			err, f := s.Apply(tx, e)
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

	tx.Commit()
	return &AddLocationsResp{
		Locs: results,
	}, nil
}
