package inventory

import (
	"context"
	"database/sql"
	"sdk-go/protos"
)

func (s *Service) AddProducts(ctx context.Context, req *protos.AddProductsReq) (r *protos.AddProductsResp, err error) {

	tx, err := s.db.Begin()
	if err != nil {
		return re(r, err)
	}

	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, "select seq from sqlite_sequence where name='Products'")
	var id uint64
	err = row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return re(r, err)
	}

	results := make([]uint64, len(req.Skus))
	for i, sku := range req.Skus {

		id += 1
		e := &protos.ProductAdded{
			Id:  id,
			Sku: sku,
		}

		err = s.Apply(tx, e)
		if err != nil {
			return re(r, err)
		}
		results[i] = id
	}

	tx.Commit()
	return &protos.AddProductsResp{Ids: results}, nil
}
