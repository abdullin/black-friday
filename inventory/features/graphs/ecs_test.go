package graphs

import (
	"fmt"
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

				actual := s.count(e.loc)
				if actual != e.qty {

					t.Errorf("Expected %d got %d on %v. \n%q", c.expect, actual, e.loc, s.ToTestString())
				}
			}

		})
	}

}
