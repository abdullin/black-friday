package stock

import (
	"black-friday/inventory/api"
	"black-friday/inventory/app"
)

func Query(ctx *app.Context, req *api.GetLocInventoryReq) (r *api.GetLocInventoryResp, err error) {

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
SELECT I.Product, SUM(I.OnHand) FROM cte_Locations AS C
JOIN Inventory AS I ON I.Location=C.Id
GROUP BY I.Product`, req.Location)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []*api.GetLocInventoryResp_Item
	for rows.Next() {
		var product uint64
		var onHand int64

		err := rows.Scan(&product, &onHand)
		if err != nil {
			return nil, err
		}

		items = append(items, &api.GetLocInventoryResp_Item{
			Product: product,
			OnHand:  onHand,
		})
	}

	rep := &api.GetLocInventoryResp{Items: items}

	return rep, nil

}
