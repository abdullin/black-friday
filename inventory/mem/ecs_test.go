package mem

import (
	"fmt"
	"math/rand"
	"testing"
)

type dat struct {
	path []int32
	qty  int32
}
type assert struct {
	qty int32
	loc int32
}

// getLocationPath gets a random but deterministic sequence
// of locations
func getLocationPath(i int32) []int32 {
	// 0 - root
	//1-10 - warehouse
	// 11-100 - shelf in a warehouse
	// 101-1000 - bin
	return []int32{
		0,
		i%9 + 1,
		i%89 + 11,
		i%889 + 111,
	}

}

func BenchmarkProductStock_Update(b *testing.B) {
	m := New()

	var i int32

	const MAX_LOCS int32 = 10000
	const MAX_PRODUCTS int32 = 1000

	// root = 11 .... 10

	for i = 0; i < MAX_LOCS; i++ {
		s := m.GetStock(int32(i % MAX_PRODUCTS))
		locs := getLocationPath(i % MAX_LOCS)

		s.Ensure(locs, 1000, 0)
	}
	for i := 0; i < b.N; i++ {
		loc := rand.Int31n(MAX_LOCS)
		product := loc % MAX_PRODUCTS

		locs := getLocationPath(loc)
		m.GetStock(product).Update(locs[len(locs)-1], 1, -1)
	}

}

func TestSomething(t *testing.T) {

	cases := []struct {
		given  []dat
		when   dat
		expect []assert
	}{
		{
			when: dat{path: []int32{0}, qty: 2},
			expect: []assert{
				{loc: 0, qty: 2},
				{loc: 100, qty: 0},
			},
		},
		{
			when: dat{path: []int32{0, 100}, qty: 2},
			expect: []assert{
				{loc: 100, qty: 2},
				{loc: 0, qty: 2},
			},
		},
		{
			given: []dat{{path: []int32{0, 100}, qty: 2}},
			when:  dat{path: []int32{0, 200}, qty: 5},
			expect: []assert{
				{loc: 0, qty: 7},
				{loc: 100, qty: 2},
				{loc: 200, qty: 5},
			},
		},
		{
			given: []dat{{path: []int32{0, 100}, qty: 2}},
			when:  dat{path: []int32{0, 100}, qty: 5},
			expect: []assert{
				{loc: 100, qty: 7},
				{loc: 0, qty: 7},
			},
		},
		{
			given: []dat{{path: []int32{0, 100}, qty: 2}},
			when:  dat{path: []int32{0, 100, 200}, qty: 5},
			expect: []assert{
				{loc: 200, qty: 5},
				{loc: 100, qty: 7},
				{loc: 0, qty: 7},
			},
		},
	}

	for i, c := range cases {
		name := fmt.Sprintf("case %d", i)
		t.Run(name, func(t *testing.T) {
			s := create()

			for _, d := range c.given {
				s.Ensure(d.path, d.qty, 0)
			}
			s.Ensure(c.when.path, c.when.qty, 0)

			for _, e := range c.expect {
				// count

				actual := s.count(e.loc)
				if actual != e.qty {

					t.Errorf("Expected %d got %d on %v. \n%q", c.expect, actual, e.loc, s.ToTestString())
				}
			}

		})
	}

}
