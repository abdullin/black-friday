package locations

import (
	"black-friday/inventory/api"
	"black-friday/inventory/app"
	"black-friday/inventory/db"
	"database/sql"
)

func List(ctx *app.Context, req *api.ListLocationsReq) (*api.ListLocationsResp, error) {

	rows, err := ctx.QueryHack(`
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
`, db.ZeroToNil(req.Location))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*api.ListLocationsResp_Loc

	lookup := make(map[uint64]*api.ListLocationsResp_Loc)

	for rows.Next() {
		var id uint64
		var parent sql.NullInt64
		var name string
		err := rows.Scan(&id, &parent, &name)
		if err != nil {
			return nil, err
		}

		loc := &api.ListLocationsResp_Loc{
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
		return nil, api.ErrNotFound
	}
	return &api.ListLocationsResp{Locs: results}, nil
}
