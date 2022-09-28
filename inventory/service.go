package inventory

import (
	c "context"
	. "sdk-go/protos"
)

type Loc struct {
	Id   uint64
	Name string
}

type product struct {
	name     string
	quantity map[uint64]int64
}

type Service struct {
	UnimplementedInventoryServiceServer

	locs     []*Loc
	products map[uint64]*product

	loc_counter  uint64
	prod_counter uint64
}

func (s *Service) AddLocation(ctx c.Context, req *AddLocationReq) (*AddLocationResp, error) {

	s.loc_counter += 1
	id := s.loc_counter
	s.locs = append(s.locs, &Loc{
		Id:   id,
		Name: req.Name,
	})

	return &AddLocationResp{Id: id}, nil
}

func (s *Service) AddProduct(ctx c.Context, req *AddProductReq) (*AddProductResp, error) {

	s.prod_counter += 1

	s.products[s.prod_counter] = &product{
		name:     req.Name,
		quantity: map[uint64]int64{},
	}

	return &AddProductResp{
		Id: s.prod_counter,
	}, nil
}

func (s *Service) UpdateQty(ctx c.Context, req *UpdateQtyReq) (*UpdateQtyResp, error) {
	//TODO implement me

	prod := s.products[req.Product]

	var current int64
	if qty, ok := prod.quantity[req.Location]; ok {
		current = qty
	}
	// TODO: handle negatives!
	prod.quantity[req.Location] = current + req.Quantity

	return &UpdateQtyResp{
		Total: current + req.Quantity,
	}, nil
}

func (s *Service) GetInventory(c.Context, *GetInventoryReq) (*GetInventoryResp, error) {

	rep := &GetInventoryResp{Items: nil}

	return rep, nil

}

func NewService() InventoryServiceServer {
	return &Service{
		products: map[uint64]*product{},
	}
}
