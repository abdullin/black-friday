package graphs

import (
	"black-friday/env/uid"
	"black-friday/inventory/api"
	"fmt"
	"google.golang.org/protobuf/proto"
)

var World = NewMem()

type Mem struct {
	stocks map[int32]*ProductStock
	locs   map[int32][]int32
	SKUs   map[string]int32
}

func NewMem() *Mem {
	return &Mem{
		stocks: map[int32]*ProductStock{},
		locs:   map[int32][]int32{0: {-1}},
		SKUs:   map[string]int32{},
	}
}

func (m *Mem) GetStock(i int32) *ProductStock {
	stock, found := m.stocks[i]
	if found {
		return stock
	}

	stock = create()
	m.stocks[i] = stock
	return stock
}

type ProductStock struct {
	locs []int32
	// root has parent index of -1
	parentIdx []int16

	reserved []int32
	onHand   []int32
}

func create() *ProductStock {
	return &ProductStock{
		locs:      []int32{0},
		parentIdx: []int16{-1},
		reserved:  []int32{0},
		onHand:    []int32{0},
	}
}

func (m *Mem) Apply(e proto.Message) {
	switch t := e.(type) {
	case *api.ProductAdded:
		m.SKUs[t.Sku] = uid.P32(t.Uid)
	case *api.InventoryUpdated:
		lid := int32(uid.Parse(t.Location))
		pid := int32(uid.Parse(t.Product))
		locs := m.locs[lid]
		m.GetStock(pid).Ensure(locs, int32(t.OnHandChange), 0)
	case *api.Reserved:
		for _, i := range t.Items {
			pid := uid.P32(i.Product)
			lid := uid.P32(i.Location)
			m.GetStock(pid).Update(lid, 0, int32(i.Quantity))
		}
	case *api.Fulfilled:
		// WHERE was the reservation?
		for _, i := range t.Items {
			pid := uid.P32(i.Product)
			lid := uid.P32(i.Location)
			m.GetStock(pid).Update(lid, int32(-i.Removed), 0)
		}

		for _, i := range t.Reserved {
			pid := uid.P32(i.Product)
			lid := uid.P32(i.Location)
			m.GetStock(pid).Update(lid, 0, int32(-i.Quantity))
		}
	case *api.LocationAdded:
		parent := uid.P32(t.Parent)
		self := uid.P32(t.Uid)

		if parent == 0 {
			m.locs[self] = []int32{0}
		} else {
			// attach to existing parent
			ancestors := m.locs[parent]
			heritage := make([]int32, len(ancestors), len(ancestors)+1)
			copy(heritage, ancestors)
			heritage = append(heritage, parent)
			m.locs[self] = heritage
		}
	}
}

func (s *ProductStock) Clone() *ProductStock {

	reserved := make([]int32, len(s.reserved))
	copy(reserved, s.reserved)
	onHand := make([]int32, len(s.onHand))
	copy(onHand, s.onHand)
	return &ProductStock{
		locs:      s.locs,
		parentIdx: s.parentIdx,
		reserved:  reserved,
		onHand:    onHand,
	}
}

func (s *ProductStock) count(loc int32) int32 {
	for i, l := range s.locs {
		if l == loc {
			return s.onHand[i]
		}
	}
	return 0
}

func (s *ProductStock) ToTestString() string {
	return fmt.Sprintf("locs: %v\nparents: %v\nonHand: %v", s.locs, s.parentIdx, s.onHand)
}

// Ensure qty to the location path
// ensure that path exists

const MAX_DEPTH = 8

func (s *ProductStock) Ensure(branch []int32, qty int32, reserve int32) {

	var bi, ti int16
	locs := s.locs

	path := make([]int16, 0, MAX_DEPTH)

	// first of all, advance as deep in the locs, as possible
	// we always advance at least one step

	for {
		//fmt.Println(fmt.Sprintf("bi %d ti %d\n", bi, ti))
		// check if we are still in the locs
		if branch[bi] == locs[ti] {

			path = append(path, ti)

			bi += 1

			// are we at the branch end?
			if bi >= int16(len(branch)) {
				break
			}
		}

		ti += 1

		// are we at the locs end?
		if ti >= int16(len(locs)) {
			break
		}
	}

	parentIdx := path[len(path)-1]
	for _, b := range branch[bi:] {
		s.parentIdx = append(s.parentIdx, parentIdx)
		s.reserved = append(s.reserved, reserve)
		s.locs = append(s.locs, b)
		s.onHand = append(s.onHand, qty)

		parentIdx = int16(len(s.locs) - 1)
	}

	// now increment the quantity

	for _, i := range path {
		s.onHand[i] += qty
		s.reserved[i] += reserve
	}

	if len(path) != len(locs) {
		panic("Ensure fail")
	}

}

func (s *ProductStock) IsValid() bool {
	// simple full tree scan
	for i := 0; i < len(s.onHand); i++ {
		if s.onHand[i] < s.reserved[i] {
			return false
		}

	}
	return true
}

// update loc hierarchy in forward
// pass (chance to reuse cache lines)
func (s *ProductStock) Update(loc int32, qty int32, reserve int32) (outQty, outReserve int32) {
	count := int16(len(s.locs))
	for i := count - 1; i >= 0; i-- {
		if s.locs[i] != loc {
			continue
		}
		outQty = s.onHand[i] + qty
		outReserve = s.reserved[i] + reserve
		s.onHand[i] = outQty
		s.reserved[i] = outReserve

		i = s.parentIdx[i]
	}
	return
}
