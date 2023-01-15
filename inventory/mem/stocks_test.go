package mem

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkStock_Update(b *testing.B) {
	m := make(map[int32][]Line)

	var i int32

	const MAX_LOCS int32 = 10000
	const MAX_PRODUCTS int32 = 1000

	// root = 11 .... 10

	for i = 0; i < MAX_LOCS; i++ {
		s := i % MAX_PRODUCTS

		locs := getLocationPath(i % MAX_LOCS)

		empty := Line{parentIdx: -1}

		lines := Ensure([]Line{empty}, locs, 1000, 0)
		m[s] = lines

	}
	for i := 0; i < b.N; i++ {
		loc := rand.Int31n(MAX_LOCS)
		product := loc % MAX_PRODUCTS

		locs := getLocationPath(loc)
		stock := m[product]
		Update(stock, locs[len(locs)-1], 1, -1)
	}

}

func TestStockLogic(t *testing.T) {

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

			locs := []Line{Line{parentIdx: -1}}

			for _, d := range c.given {
				locs = Ensure(locs, d.path, d.qty, 0)
			}
			locs = Ensure(locs, c.when.path, c.when.qty, 0)

			for _, e := range c.expect {

				var actual int32
				// counting
				for _, l := range locs {
					if l.loc == e.loc {
						actual = l.onHand
						break
					}
				}

				if actual != e.qty {

					t.Errorf("Expected %d got %d on %v. \n%q", c.expect, actual, e.loc, ToTestString(locs))
				}
			}

		})
	}

}
