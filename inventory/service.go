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
	store *Store
	UnimplementedInventoryServiceServer
}

func (s *Service) ListLocations(ctx c.Context, req *ListLocationsReq) (*ListLocationsResp, error) {

	results := make([]*ListLocationsResp_Loc, len(s.store.locs))
	for i, l := range s.store.locs {
		results[i] = &ListLocationsResp_Loc{
			Id:   l.Id,
			Name: l.Name,
		}

	}
	return &ListLocationsResp{Locs: results}, nil
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

func (s *Service) AddLocations(_ c.Context, req *AddLocationsReq) (*AddLocationsResp, error) {

	results := make([]uint64, len(req.Names))
	for i, name := range req.Names {
		var id = s.store.loc_counter + 1

		e := &LocationAdded{
			Name: name,
			Id:   id,
		}
		results[i] = id

		s.store.Apply(e)
	}

	return &AddLocationsResp{Ids: results}, nil
}

func (s *Service) AddProducts(ctx c.Context, req *AddProductsReq) (*AddProductsResp, error) {

	results := make([]uint64, len(req.Skus))
	for i, sku := range req.Skus {
		if _, found := s.store.products_index[sku]; found {
			return nil, status.Errorf(codes.AlreadyExists, "SKU %s already exists", sku)
		}

		id := s.store.prod_counter + 1

		e := &ProductAdded{
			Id:  id,
			Sku: sku,
		}
		s.store.Apply(e)

		results[i] = id

	}

	return &AddProductsResp{Ids: results}, nil
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
		return nil, status.Errorf(codes.FailedPrecondition, "Can't be negative!")
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
