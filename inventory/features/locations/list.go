package locations

import (
	"black-friday/env/uid"
	"black-friday/fx"
	"black-friday/inventory/api"
	"google.golang.org/grpc/status"
)

func List(ctx fx.Tx, req *api.ListLocationsReq) (*api.ListLocationsResp, *status.Status) {

	id := uid.Parse(req.Location)
	rows, err := ctx.QueryHack(`
WITH RECURSIVE cte_Locations(Id, Parent, Name) AS (
	SELECT l.Id, l.Parent, l.Name
	FROM Locations l
	WHERE CASE
		WHEN ? = 0
		THEN l.Parent = 0 and l.Id != 0
		ELSE l.Id = ?
	END 

	UNION ALL

	SELECT l.Id, l.Parent, l.Name
	FROM Locations l
	JOIN cte_Locations c ON c.Id = l.Parent
)
SELECT * FROM cte_Locations
`, id, id)
	if err != nil {
		return nil, status.Convert(err)
	}
	defer rows.Close()

	var results []*api.ListLocationsResp_Loc

	lookup := make(map[int64]*api.ListLocationsResp_Loc)

	for rows.Next() {
		var id int64
		var parent int64
		var name string
		err := rows.Scan(&id, &parent, &name)
		if err != nil {
			return nil, status.Convert(err)
		}

		loc := &api.ListLocationsResp_Loc{
			Name:    name,
			Uid:     uid.Str(id),
			Parent:  uid.Str(parent),
			Chidren: nil,
		}
		lookup[id] = loc
		if parent, found := lookup[parent]; found {
			parent.Chidren = append(parent.Chidren, loc)
		} else {
			results = append(results, loc)
		}
	}
	// we should at list get one location
	if len(results) == 0 {
		return nil, api.ErrLocationNotFound
	}
	return &api.ListLocationsResp{Locs: results}, nil
}
