package inventory

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	. "sdk-go/protos"
)

func (s *Service) UpdateInventory(ctx context.Context, req *UpdateInventoryReq) (r *UpdateInventoryResp, err error) {

	tx, err := s.db.Begin()
	if err != nil {
		return re(r, err)
	}

	defer tx.Rollback()

	row := tx.QueryRowContext(ctx,
		"SELECT OnHand FROM Inventory WHERE Location=? AND Product=?",
		req.Location,
		req.Product)

	var onHand int64

	err = row.Scan(&onHand)
	if err != nil && err != sql.ErrNoRows {
		return re(r, err)
	}

	onHand += req.OnHandChange

	if onHand < 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "Can't be negative!")
	}

	e := &InventoryUpdated{
		Location:     req.Location,
		Product:      req.Product,
		OnHandChange: req.OnHandChange,
		OnHand:       onHand,
	}

	err = s.Apply(tx, e)
	if err != nil {
		return re(r, err)
	}

	tx.Commit()

	return &UpdateInventoryResp{
		OnHand: e.OnHand,
	}, nil
}

func (s *Service) GetInventory(c context.Context, req *GetInventoryReq) (r *GetInventoryResp, err error) {

	tx, err := s.db.Begin()
	if err != nil {
		return re(r, err)
	}

	defer tx.Rollback()

	rows, err := tx.QueryContext(c, "SELECT Product, OnHand FROM Inventory WHERE Location=?", req.Location)
	if err != nil {
		return re(r, err)
	}

	var items []*GetInventoryResp_Item
	for rows.Next() {
		var product uint64
		var onHand int64

		err := rows.Scan(&product, &onHand)
		if err != nil {
			return re(r, err)
		}

		items = append(items, &GetInventoryResp_Item{
			Product: product,
			OnHand:  onHand,
		})
	}

	rep := &GetInventoryResp{Items: items}

	tx.Commit()

	return rep, nil

}
