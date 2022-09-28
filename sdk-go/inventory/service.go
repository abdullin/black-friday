package inventory

import (
	c "context"
	. "sdk-go/protos"
)

type Loc struct {
	Id   uint64
	Name string
}

type Service struct {
	UnimplementedInventoryServiceServer

	locs []*Loc

	counter uint64
}

func (s *Service) AddLocation(ctx c.Context, req *AddLocationRequest) (*AddLocationResponse, error) {

	s.counter += 1
	id := s.counter
	s.locs = append(s.locs, &Loc{
		Id:   id,
		Name: req.Name,
	})

	return &AddLocationResponse{Id: id}, nil
}

func (s *Service) AddProduct(ctx c.Context, req *AddProductRequest) (*AddProductResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ChangeQuantity(ctx c.Context, req *ChangeQuantityRequest) (*ChangeQuantityResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ListLocation(c.Context, *ListLocationRequest) (*ListLocationResponse, error) {

	rep := &ListLocationResponse{Items: nil}

	return rep, nil

}

func NewService() InventoryServiceServer {
	return &Service{}
}
