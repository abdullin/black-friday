package tests

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"sdk-go/inventory"
	. "sdk-go/protos"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func TestSomething(t *testing.T) {
	check := func(err error) {
		if err != nil {
			panic(err)
			//t.Fatal(err)
		}
	}

	db, err := sql.Open("sqlite3", ":memory:")
	check(err)

	defer db.Close()

	check(inventory.CreateSchema(db))

	s := inventory.NewService(db)

	ctx := context.Background()

	prods, err := s.AddProducts(ctx, &AddProductsReq{Skus: []string{"One", "Two"}})
	check(err)
	locs, err := s.AddLocations(ctx, &AddLocationsReq{Names: []string{"Cola", "Fanta"}})
	check(err)

	q, err := s.UpdateQty(ctx, &UpdateQtyReq{
		Location: locs.Ids[0],
		Product:  prods.Ids[0],
		Quantity: 1,
	})

	check(err)

	t.Log(q.Total)

	_, _ = s.Reset(nil, nil)

}

func BenchmarkCreation(b *testing.B) {

	ctx := context.Background()
	check := func(err error) {
		if err != nil {
			b.Fatal(err)
		}
	}

	dir := b.TempDir()

	db, err := sql.Open("sqlite3", path.Join(dir, "test.sqlite"))
	check(err)

	defer db.Close()

	check(inventory.CreateSchema(db))
	s := inventory.NewService(db)

	start := time.Now()

	for i := 0; i < b.N; i++ {

		sku := fmt.Sprintf("SKU_%d", i)
		shelf := fmt.Sprintf("SHELF_%d", i)

		prods, err := s.AddProducts(ctx, &AddProductsReq{Skus: []string{sku}})
		check(err)
		shelves, err := s.AddLocations(ctx, &AddLocationsReq{Names: []string{shelf}})
		check(err)

		_, err = s.UpdateQty(ctx, &UpdateQtyReq{
			Location: shelves.Ids[0],
			Product:  prods.Ids[0],
			Quantity: 1,
		})
		check(err)
	}

	elapsed := time.Since(start)
	frequency := float64(b.N) / elapsed.Seconds()

	b.Logf("Size %d in %f, frequency: %f\n", b.N, elapsed.Seconds(), frequency)
	check(err)
}
