package stock

import (
	"black-friday/fx"
	"black-friday/inventory/api"
)

func Query(ctx fx.Tx, req *api.GetLocInventoryReq) (r *api.GetLocInventoryResp, err error) {

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
SELECT I.Product, SUM(I.OnHand), IFNULL(SUM(R.Quantity),0) 
FROM cte_Locations AS C
JOIN Inventory AS I ON I.Location=C.Id
LEFT JOIN Reserves AS R ON R.Product=I.Product AND R.Location=I.Location
GROUP BY I.Product`, req.Location)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []*api.GetLocInventoryResp_Item
	for rows.Next() {
		var product int64
		var onHand int64
		var reserved int64

		err := rows.Scan(&product, &onHand, &reserved)
		if err != nil {
			return nil, err
		}

		items = append(items, &api.GetLocInventoryResp_Item{
			Product:   product,
			OnHand:    onHand,
			Available: onHand - reserved,
		})
	}

	rep := &api.GetLocInventoryResp{Items: items}

	return rep, nil

}
