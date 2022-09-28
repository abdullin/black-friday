package inventory

import (
	c "context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
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

	store *Store
}

type Store struct {
	locs           []*Loc
	products       map[uint64]*product
	products_index map[string]uint64

	loc_counter  uint64
	prod_counter uint64
}

func (s *Store) Apply(e proto.Message) {

	switch t := e.(type) {
	case *LocationAdded:
		s.locs = append(s.locs, &Loc{
			Id:   t.Id,
			Name: t.Name,
		})
		s.loc_counter = t.Id

	case *ProductAdded:
		s.products[t.Id] = &product{
			name:     t.Sku,
			quantity: map[uint64]int64{},
		}
		s.products_index[t.Sku] = t.Id
		s.prod_counter = t.Id
	case *QuantityUpdated:
		s.products[t.Product].quantity[t.Location] = t.Total

	default:
		panic("UNKNOWN EVENT")

	}

}

func (s *Service) AddLocation(ctx c.Context, req *AddLocationReq) (*AddLocationResp, error) {

	e := &LocationAdded{
		Name: req.Name,
		Id:   s.store.loc_counter + 1,
	}

	s.store.Apply(e)

	return &AddLocationResp{Id: e.Id}, nil
}

func (s *Service) AddProduct(ctx c.Context, req *AddProductReq) (*AddProductResp, error) {

	if _, found := s.store.products_index[req.Sku]; found {
		return nil, status.Errorf(codes.AlreadyExists, "SKU %s already exists", req.Sku)
	}

	e := &ProductAdded{
		Id:  s.store.prod_counter + 1,
		Sku: req.Sku,
	}
	s.store.Apply(e)

	return &AddProductResp{
		Id: e.Id,
	}, nil
}

func (s *Service) UpdateQty(ctx c.Context, req *UpdateQtyReq) (*UpdateQtyResp, error) {
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
		return nil, status.Errorf(codes.InvalidArgument, "Can't be negative!")
	}

	e := &QuantityUpdated{
		Location: req.Location,
		Product:  req.Product,
		Quantity: req.Quantity,
		Total:    total,
	}

	s.store.Apply(e)

	return &UpdateQtyResp{
		Total: e.Total,
	}, nil
}

func (s *Service) GetInventory(c c.Context, r *GetInventoryReq) (*GetInventoryResp, error) {

	var items []*GetInventoryResp_Item

	for id, p := range s.store.products {
		if qty, found := p.quantity[r.Location]; found && qty != 0 {
			items = append(items, &GetInventoryResp_Item{
				Product:  id,
				Quantity: qty,
			})
		}
	}

	rep := &GetInventoryResp{Items: items}

	return rep, nil

}

func NewService() InventoryServiceServer {

	s := &Store{

		products:       map[uint64]*product{},
		products_index: map[string]uint64{},
	}
	return &Service{
		store: s,
	}
}
