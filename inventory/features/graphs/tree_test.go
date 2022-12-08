package graphs

import "testing"

func TestWalk(t *testing.T) {
	data := []struct {
		n    *Node
		ok   bool
		name string
	}{

		{
			n:    &Node{},
			ok:   true,
			name: "simple empty",
		},
		{
			n: &Node{
				OnHand:   10,
				Reserved: 10,
			},
			ok:   true,
			name: "simple filled",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			_, _, ok := Walk(d.n)
			if ok != d.ok {
				t.Fatalf("expected %v actual %v", d.ok, ok)
			}
		})
	}
}
