package inventory

import (
	"context"
	"database/sql"
	"sdk-go/fail"
	. "sdk-go/protos"
	"sdk-go/stat"
)

func (s *Service) ListLocations(ctx context.Context, req *ListLocationsReq) (r *ListLocationsResp, e error) {

	tx := s.GetTx(ctx)

	var loc any = nil
	if req.Location != 0 {
		loc = req.Location
	}

	rows, err := tx.tx.QueryContext(ctx, `
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
`, loc)
	if err != nil {
		return re(r, err)
	}
	defer rows.Close()

	var results []*ListLocationsResp_Loc

	for rows.Next() {
		var id uint64
		var parent sql.NullInt64
		var name string
		err := rows.Scan(&id, &parent, &name)
		if err != nil {
			return re(r, err)
		}
		results = append(results, &ListLocationsResp_Loc{
			Name:    name,
			Id:      id,
			Parent:  uint64(parent.Int64),
			Chidren: nil,
		})
	}
	// we should at list get one location
	if len(results) == 0 {
		return nil, stat.NotFound
	}
	return &ListLocationsResp{Locs: results}, nil
}

func (s *Service) AddLocations(ctx context.Context, req *AddLocationsReq) (r *AddLocationsResp, e error) {

	tx := s.GetTx(ctx)

	id := tx.GetSeq("Locations")

	var addLoc func(parent uint64, ls []*AddLocationsReq_Loc) ([]*AddLocationsResp_Loc, error)

	addLoc = func(parent uint64, ls []*AddLocationsReq_Loc) ([]*AddLocationsResp_Loc, error) {

		var r []*AddLocationsResp_Loc

		for _, l := range ls {
			if l.Name == "" {
				return nil, stat.ArgNil("Name")
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

			err, f := tx.Apply(e)
			switch f {
			case fail.OK:
			case fail.ConstraintUnique:
				return nil, stat.DuplicateName
			case fail.ConstraintForeign:
				return nil, stat.NotFound
			default:
				return nil, stat.Internal(err, f)
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
