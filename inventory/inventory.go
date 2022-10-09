package inventory

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sdk-go/protos"
)

func (s *Service) UpdateQty(ctx context.Context, req *protos.UpdateQtyReq) (*protos.UpdateQtyResp, error) {
	//TODO implement me

	prod := s.store.products[req.Product]

	if prod == nil {
		log.Panicln("NIIIL for product ", req.Product)
	}

	var current int64
	if qty, ok := prod.quantity[req.Location]; ok {
		current = qty
	}

	total := current + req.Quantity

	if total < 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "Can't be negative!")
	}

	e := &protos.QuantityUpdated{
		Location: req.Location,
		Product:  req.Product,
		Quantity: req.Quantity,
		Total:    total,
	}

	s.store.Apply(e)

	return &protos.UpdateQtyResp{
		Total: e.Total,
	}, nil
}

func (s *Service) GetInventory(c context.Context, r *protos.GetInventoryReq) (*protos.GetInventoryResp, error) {

	var items []*protos.GetInventoryResp_Item

	for id, p := range s.store.products {
		if qty, found := p.quantity[r.Location]; found && qty != 0 {
			items = append(items, &protos.GetInventoryResp_Item{
				Product:  id,
				Quantity: qty,
			})
		}
	}

	rep := &protos.GetInventoryResp{Items: items}

	return rep, nil

}
