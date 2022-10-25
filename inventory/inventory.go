package inventory

import (
	. "black-friday/api"
	"black-friday/fail"
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) UpdateInventory(ctx context.Context, req *UpdateInventoryReq) (r *UpdateInventoryResp, err error) {
	tx := s.GetTx(ctx)

	onHand, err := tx.QueryInt64("SELECT OnHand FROM Inventory WHERE Location=? AND Product=?",
		req.Location,
		req.Product)

	if err != sql.ErrNoRows {
		return nil, err
	}

	onHand += req.OnHandChange

	if onHand < 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "OnHand can't go negative!")
	}

	e := &InventoryUpdated{
		Location:     req.Location,
		Product:      req.Product,
		OnHandChange: req.OnHandChange,
		OnHand:       onHand,
	}

	err, f := s.Apply(tx, e)
	switch f {
	case fail.None:
	default:
		return nil, ErrInternal(err, f)
	}

	tx.Commit()

	return &UpdateInventoryResp{OnHand: e.OnHand}, nil
}

func (s *Service) GetLocInventory(ctx context.Context, req *GetLocInventoryReq) (r *GetLocInventoryResp, err error) {
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
SELECT I.Product, SUM(I.OnHand) FROM cte_Locations AS C
JOIN Inventory AS I ON I.Location=C.Id
GROUP BY I.Product`, req.Location)
	if err != nil {
		return nil, err
	}

	var items []*GetLocInventoryResp_Item
	for rows.Next() {
		var product uint64
		var onHand int64

		err := rows.Scan(&product, &onHand)
		if err != nil {
			return nil, err
		}

		items = append(items, &GetLocInventoryResp_Item{
			Product: product,
			OnHand:  onHand,
		})
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	rep := &GetLocInventoryResp{Items: items}

	tx.Commit()

	return rep, nil

}
