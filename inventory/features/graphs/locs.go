package graphs

import (
	"black-friday/fx"
	"database/sql"
)

// map - loc/id -> leaf node -> chain up to root.

var Cache map[int64]int64

func LoadLocTree(ctx fx.Tx) (map[int64]int64, error) {
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
	Cache = m
	return m, nil
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

func Resolves(locs map[int64]int64, onHand, reserved []Stock) bool {
	// cached tree
	var stocks = make(map[int64]struct{ hand, reserved int64 })

	for _, a := range onHand {
		current := a.Loc

		// walk this up, summing values in stocks
		for {
			stock := stocks[current]
			stock.hand += a.Qty
			stocks[current] = stock

			if current == 0 {
				break
			}
			current = locs[current]
		}
	}

	for _, a := range reserved {

		current := a.Loc
		// walk this up, summing values in stocks
		for {
			stock := stocks[current]
			stock.reserved += a.Qty

			if stock.reserved > stock.hand {
				return false
			}

			stocks[current] = stock

			if current == 0 {
				break
			}
			current = locs[current]
		}
	}

	return true
}
