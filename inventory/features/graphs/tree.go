package graphs

import (
	"black-friday/fx"
	"fmt"
	"strings"
)

type Node struct {
	Location int64
	Parent   int64
	OnHand   int64
	Reserved int64

	Children []*Node
}

func Modify(n *Node, loc, onHand, reserve int64) (int64, int64, bool) {

	if n.Location == loc {
		n.OnHand += onHand
		n.Reserved += reserve
		return n.OnHand, n.Reserved, true
	}

	for _, c := range n.Children {
		a, b, c := Modify(c, loc, onHand, reserve)
		if c {
			return a, b, c
		}
	}
	return 0, 0, false

}

func Print(n *Node, level int) {
	fmt.Printf("%s %d (%d, %d)\n", strings.Repeat("  ", level), n.Location, n.OnHand, n.Reserved)
	for _, i := range n.Children {
		Print(i, level+1)
	}

}

func Walk(root *Node) (onHand, reserved int64, ok bool) {
	// we need to walk down and sum up all quantities and availabilities

	onHand = root.OnHand
	reserved = root.Reserved
	ok = true

	for _, x := range root.Children {
		a, b, c := Walk(x)
		onHand += a
		reserved += b

		ok = ok && c
	}
	ok = ok && (reserved <= onHand)
	return
}

// LoadProductTree goes through a list of all locations that contain
// a product quantity or a reservation. Empty Nodes are pruned
func LoadProductTree(ctx fx.Tx, product int64) (*Node, error) {
	query := `


SELECT 
	T.Location AS Location, 
	LL.Parent AS Parent, 
	ifnull(SUM(I.OnHand),0) AS OnHand, 
	ifnull(SUM(R.Quantity),0) AS Reserved	
FROM (
	WITH RECURSIVE TREE (Location) AS (
	-- SEED
	SELECT DISTINCT Location FROM INVENTORY
	WHERE PRODUCT=?
	
	UNION ALL
	
	SELECT l.Parent
	FROM Locations l
	JOIN TREE T ON T.Location = l.Id
	WHERE T.Location != 0
)
SELECT DISTINCT Location FROM TREE) AS T
LEFT JOIN INVENTORY I ON T.Location=I.Location AND I.Product=?
LEFT JOIN Reserves R ON T.Location=R.Location AND R.Product=?
LEFT JOIN Locations LL on T.Location=LL.Id

GROUP BY T.Location
`

	ctx.Trace().Begin("Load Product")
	defer ctx.Trace().End()
	rows, err := ctx.QueryHack(query, product, product, product)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lookup := make(map[int64]*Node)
	for rows.Next() {
		n := &Node{}
		err := rows.Scan(&n.Location, &n.Parent, &n.OnHand, &n.Reserved)
		if err != nil {
			return nil, err
		}
		lookup[n.Location] = n
	}
	var root *Node

	for id, node := range lookup {
		if id == 0 {
			root = node
		} else {
			parent := lookup[node.Parent]
			parent.Children = append(parent.Children, node)
		}
	}
	if nil == root {
		// no inventory at all
		return &Node{}, nil
	}

	return root, nil
}
