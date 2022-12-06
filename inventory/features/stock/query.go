package stock

import (
	"black-friday/env/uid"
	"black-friday/fx"
	"black-friday/inventory/api"
	"google.golang.org/grpc/status"
)

type Quantity struct {
	Product  int64
	Quantity int64
}

func walkQuantitiesDown(ctx fx.Tx, location int64) ([]Quantity, error) {
	rows, err := ctx.QueryHack(`
WITH RECURSIVE cte_Locations(Id, Parent, Name) AS (
	SELECT l.Id, l.Parent, l.Name
	FROM Locations l
	WHERE CASE
		WHEN ? = 0 
		THEN l.Parent=0 and l.Id != 0
		ELSE l.id = ?
	END

	UNION ALL

	SELECT l.Id, l.Parent, l.Name
	FROM Locations l
	JOIN cte_Locations c ON c.Id = l.Parent
)
SELECT I.Product AS Product, SUM(I.OnHand) AS OnHand
FROM cte_Locations AS C
JOIN Inventory AS I ON I.Location=C.Id
GROUP BY Product`, location, location)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var quantities []Quantity

	for rows.Next() {
		var product int64
		var onHand int64

		err := rows.Scan(&product, &onHand)
		if err != nil {
			return nil, err
		}
		quantities = append(quantities, Quantity{
			Product:  product,
			Quantity: onHand,
		})

	}

	return quantities, nil

}

func walkReservationsDown(ctx fx.Tx, location int64) (map[int64]int64, error) {
	rows, err := ctx.QueryHack(`
WITH RECURSIVE cte_Locations(Id, Parent, Name) AS (
	SELECT l.Id, l.Parent, l.Name
	FROM Locations l
	WHERE CASE
		WHEN ? = 0 
		THEN l.Parent=0 and l.Id != 0
		ELSE l.id = ?
	END

	UNION ALL

	SELECT l.Id, l.Parent, l.Name
	FROM Locations l
	JOIN cte_Locations c ON c.Id = l.Parent
)
SELECT R.Product, SUM(R.Quantity)
FROM cte_Locations AS C
JOIN Reserves AS R ON R.Location=C.Id
GROUP BY R.Product`, location, location)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	quantities := make(map[int64]int64)

	for rows.Next() {
		var product int64
		var onHand int64

		err := rows.Scan(&product, &onHand)
		if err != nil {
			return nil, err
		}
		quantities[product] = onHand

	}

	return quantities, nil

}

func Query(ctx fx.Tx, req *api.GetLocInventoryReq) (*api.GetLocInventoryResp, *status.Status) {

	loc := uid.Parse(req.Location)

	quantities, err := walkQuantitiesDown(ctx, loc)
	if err != nil {
		return nil, status.Convert(err)
	}
	reserves, err := walkReservationsDown(ctx, loc)
	if err != nil {
		return nil, status.Convert(err)
	}

	var items []*api.GetLocInventoryResp_Item

	for _, q := range quantities {
		onHand := q.Quantity
		totalReserved, _ := reserves[q.Product]

		items = append(items, &api.GetLocInventoryResp_Item{
			Product:   uid.Str(q.Product),
			OnHand:    onHand,
			Available: onHand - totalReserved,
		})

	}

	rep := &api.GetLocInventoryResp{Items: items}

	return rep, nil

}
