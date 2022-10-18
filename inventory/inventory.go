package inventory

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sdk-go/fail"
	. "sdk-go/protos"
	"sdk-go/stat"
)

func (s *Service) UpdateInventory(ctx context.Context, req *UpdateInventoryReq) (r *UpdateInventoryResp, err error) {
	tx := s.GetTx(ctx)

	onHand, err := tx.QueryInt64("SELECT OnHand FROM Inventory WHERE Location=? AND Product=?",
		req.Location,
		req.Product)

	if err != sql.ErrNoRows {
		return re(r, err)
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

	err, f := tx.Apply(e)
	switch f {
	case fail.OK:
	default:
		return nil, stat.Internal(err, f)
	}

	tx.Commit()

	return &UpdateInventoryResp{OnHand: e.OnHand}, nil
}

func (s *Service) GetInventory(ctx context.Context, req *GetInventoryReq) (r *GetInventoryResp, err error) {

	tx := s.GetTx(ctx)

	rows, err := tx.tx.QueryContext(ctx, "SELECT Product, OnHand FROM Inventory WHERE Location=?", req.Location)
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

	if err := rows.Close(); err != nil {
		return nil, err
	}

	rep := &GetInventoryResp{Items: items}

	tx.Commit()

	return rep, nil

}
