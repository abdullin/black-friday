package inventory

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sdk-go/protos"
)

func (s *Service) UpdateQty(ctx context.Context, req *protos.UpdateQtyReq) (r *protos.UpdateQtyResp, err error) {

	tx, err := s.db.Begin()
	if err != nil {
		return re(r, err)
	}

	defer tx.Rollback()

	row := tx.QueryRowContext(ctx,
		"SELECT Quantity FROM Inventory WHERE Location=? AND Product=?",
		req.Location,
		req.Product)

	var quantity int64

	err = row.Scan(&quantity)
	if err != nil && err != sql.ErrNoRows {
		return re(r, err)
	}

	total := quantity + req.Quantity

	if total < 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "Can't be negative!")
	}

	e := &protos.QuantityUpdated{
		Location: req.Location,
		Product:  req.Product,
		Quantity: req.Quantity,
		Total:    total,
		Before:   quantity,
	}

	err = s.Apply(tx, e)
	if err != nil {
		return re(r, err)
	}

	tx.Commit()

	return &protos.UpdateQtyResp{
		Total: e.Total,
	}, nil
}

func (s *Service) GetInventory(c context.Context, req *protos.GetInventoryReq) (r *protos.GetInventoryResp, err error) {

	tx, err := s.db.Begin()
	if err != nil {
		return re(r, err)
	}

	defer tx.Rollback()

	rows, err := tx.QueryContext(c, "SELECT Product, Quantity FROM Inventory WHERE Location=?", req.Location)
	if err != nil {
		return re(r, err)
	}

	var items []*protos.GetInventoryResp_Item
	for rows.Next() {
		var product uint64
		var quantity int64

		err := rows.Scan(&product, &quantity)
		if err != nil {
			return re(r, err)
		}

		items = append(items, &protos.GetInventoryResp_Item{
			Product:  product,
			Quantity: quantity,
		})
	}

	rep := &protos.GetInventoryResp{Items: items}

	tx.Commit()

	return rep, nil

}
