package graphs

import (
	"black-friday/fx"
	"database/sql"
)

// map - loc/id -> leaf node -> chain up to root.

var Cache Tree

type Tree map[int64][]int64

func DenormaliseLocTree(in map[int64]int64) Tree {
	t := make(Tree)

	for id, parent := range in {

		var parents []int64

		for {
			parents = append(parents, parent)
			if parent == 0 {
				break
			} else {
				parent = in[parent]
			}

		}
		t[id] = parents

	}
	return t

}

func LoadLocTree(ctx fx.Tx) (Tree, error) {
	if Cache != nil {
		return Cache, nil
	}

	m := make(map[int64]int64)

	query := `SELECT Id, Parent from Locations`

	rows, err := ctx.QueryHack(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, parent int64
		err := rows.Scan(&id, &parent)
		if err != nil {
			return nil, err
		}
		m[id] = parent
	}
	Cache = DenormaliseLocTree(m)
	return Cache, nil
}

type Stock struct {
	Qty, Loc int64
}

func LoadStocks(ctx fx.Tx, product int64) (reserves []Stock, inventory []Stock, err error) {

	query := `SELECT Location, OnHand, 1 FROM Inventory WHERE Product=?
UNION
SELECT Location, SUM(Quantity), 0 FROM Reserves WHERE Product=?
GROUP BY LOcation`
	var rows *sql.Rows
	rows, err = ctx.QueryHack(query, product, product)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var s Stock
		var isInventory int

		err = rows.Scan(&s.Loc, &s.Qty, &isInventory)
		if err != nil {
			return
		}
		if isInventory == 1 {
			inventory = append(inventory, s)
		} else {
			reserves = append(reserves, s)
		}
	}
	return
}

func LoadReserves(ctx fx.Tx, product int64) ([]Stock, error) {

	ctx.Trace().Begin("Load Reserves")
	defer ctx.Trace().End()
	var res []Stock
	rows, err := ctx.QueryHack("SELECT Location, Sum(Quantity) FROM Reserves WHERE Product=? GROUP BY Location", product)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var s Stock
		err := rows.Scan(&s.Loc, &s.Qty)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func LoadInventory(ctx fx.Tx, product int64) ([]Stock, error) {

	ctx.Trace().Begin("Load Inventory")
	defer ctx.Trace().End()
	var res []Stock
	rows, err := ctx.QueryHack("SELECT Location, OnHand FROM Inventory WHERE Product=?", product)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var s Stock
		err := rows.Scan(&s.Loc, &s.Qty)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func Resolves(ctx fx.Tx, locs Tree, onHand, reserved []Stock) bool {

	ctx.Trace().Begin("Resolve")
	defer ctx.Trace().End()
	// cached tree
	var stocks = make(map[int64]struct{ hand, reserved int64 }, len(reserved)*2)

	for _, a := range onHand {

		stock := stocks[a.Loc]
		stock.hand += a.Qty
		stocks[a.Loc] = stock

		parents := locs[a.Loc]

		for _, parent := range parents {
			stock := stocks[parent]
			stock.hand += a.Qty
			stocks[parent] = stock
		}
	}

	for _, a := range reserved {
		stock := stocks[a.Loc]
		stock.reserved += a.Qty

		if stock.reserved > stock.hand {
			return false
		}

		stocks[a.Loc] = stock

		for _, parent := range locs[a.Loc] {
			stock := stocks[parent]
			stock.reserved += a.Qty

			if stock.reserved > stock.hand {
				return false
			}

			stocks[parent] = stock
		}
	}

	return true
}
