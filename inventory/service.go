package inventory

import (
	"google.golang.org/protobuf/proto"
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

func NewService() InventoryServiceServer {

	s := &Store{

		products:       map[uint64]*product{},
		products_index: map[string]uint64{},
	}
	return &Service{
		store: s,
	}
}
