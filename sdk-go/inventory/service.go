package inventory

import (
	"context"
	p "sdk-go/protos"
)

type Loc struct {
	Id   uint64
	Name string
}

type Service struct {
	p.UnimplementedInventoryServiceServer

	locs []*Loc

	counter uint64
}

func (s *Service) AddLocation(ctx context.Context, request *p.AddLocationRequest) (*p.AddLocationResponse, error) {

	s.counter += 1
	id := s.counter
	s.locs = append(s.locs, &Loc{
		Id:   id,
		Name: request.Name,
	})

	return &p.AddLocationResponse{Id: id}, nil
}

func (s *Service) AddProduct(ctx context.Context, request *p.AddProductRequest) (*p.AddProductResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ChangeQuantity(ctx context.Context, request *p.ChangeQuantityRequest) (*p.ChangeQuantityResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) mustEmbedUnimplementedInventoryServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ListLocation(context.Context, *p.ListLocationRequest) (*p.ListLocationResponse, error) {

	rep := &p.ListLocationResponse{Items: nil}

	return rep, nil

}

func NewService() p.InventoryServiceServer {
	return &Service{}
}
