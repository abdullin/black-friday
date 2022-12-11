package stress

import (
	"black-friday/env/rnd"
	"black-friday/env/uid"
	"black-friday/inventory/api"
	"container/list"
	"context"
	"fmt"
	"log"
)

type env struct {
	products   int64
	locations  int64
	warehouses int64
	inventory  []int64

	reject       int64
	sales        int64
	reservations list.List

	bins      []string
	fulfilled int64

	r *rnd.Rand

	client api.InventoryServiceClient
}

func NewEnv(client api.InventoryServiceClient) *env {
	return &env{client: client, r: rnd.New(), inventory: make([]int64, 10000, 10000)}
}

func (e *env) TryFulfull(ctx context.Context, count int) {
	for j := 0; j < count; j++ {
		if e.reservations.Len() == 0 {
			return
		}
		loc := int(e.r.Int63n(int64(e.reservations.Len())))
		n := e.reservations.Front()
		for i := 0; i < loc; i++ {
			n = n.Next()
		}
		id := n.Value.(string)
		_, err := e.client.Fulfill(ctx, &api.FulfillReq{
			Reservation: id,
		})
		e.reservations.Remove(n)
		if err != nil {
			log.Panicln(err)
		} else {
			e.fulfilled += 1
		}

	}

}

func (e *env) TrySell(ctx context.Context) {
	name := fmt.Sprintf("sale-%d", e.sales)

	c := int(e.r.Int63n(10) + 1)

	var items []*api.ReserveReq_Item

	prods := make(map[int64]struct{})

	for i := 0; i < c; i++ {

		product := e.r.Int63n(e.products) + 1
		// try selling something that is in store
		for j := product; j < product+200; j++ {
			if e.inventory[j] > 0 {
				product = j
				break
			}
		}

		if _, found := prods[product]; found {
			continue
		}
		prods[product] = struct{}{}

		items = append(items, &api.ReserveReq_Item{
			Sku:      SKU(product),
			Quantity: e.r.Int63n(5) + 1,
		})

	}

	if len(items) == 0 {
		return
	}

	e.sales += 1

	r, err := e.client.Reserve(ctx, &api.ReserveReq{
		Reservation: name,
		Items:       items,
	})

	if err == nil {
		e.reservations.PushBack(r.Reservation)

	} else {
		e.sales -= 1
		e.reject += 1
	}

}

func (e *env) AddInventory(ctx context.Context) {

	product := e.r.Int63n(e.products) + 1
	location := e.bins[product-1]

	quantity := e.r.Int63n(200) + 20

	_, err := e.client.UpdateInventory(ctx, &api.UpdateInventoryReq{
		Location:     location,
		Product:      uid.Str(product),
		OnHandChange: quantity,
	})

	e.inventory[product] += quantity

	if err != nil {
		log.Fatalln(err)
	}

	// which product?

	//
}

func (e *env) AddProducts(ctx context.Context) {
	e.products += 1

	var skus []string

	skus = append(skus, SKU(e.products))

	prod := &api.AddProductsReq{Skus: skus}

	_, err := e.client.AddProducts(ctx, prod)
	if err != nil {
		log.Panicln(err)
	}

}

func SKU(e int64) string {
	return fmt.Sprintf("product-%d", e)
}

func (e *env) AddWarehouse(ctx context.Context) {

	e.warehouses += 1

	whsName := fmt.Sprintf("WHS-%d", e.warehouses)

	whs := &api.AddLocationsReq_Loc{
		Name: whsName,
	}

	e.locations += 1

	// add rows
	for r := 0; r < 3; r++ {
		rowName := fmt.Sprintf("%s/ROW-%d", whsName, r+1)
		row := &api.AddLocationsReq_Loc{
			Name: rowName}
		whs.Locs = append(whs.Locs, row)

		e.locations += 1

		for s := 0; s < 4; s++ {
			shelfName := fmt.Sprintf("%s/SHELF-%d", rowName, s+1)
			shelf := &api.AddLocationsReq_Loc{Name: shelfName}
			row.Locs = append(row.Locs, shelf)

			e.locations += 1

			for b := 0; b < 5; b++ {
				binName := fmt.Sprintf("BIN-%d", e.locations)
				bin := &api.AddLocationsReq_Loc{Name: binName}
				shelf.Locs = append(shelf.Locs, bin)

				e.locations += 1
			}

		}

	}

	resp, err := e.client.AddLocations(ctx, &api.AddLocationsReq{
		Locs:   []*api.AddLocationsReq_Loc{whs},
		Parent: uid.Str(0),
	})
	if err != nil {
		log.Panicln(err)
	}

	for _, w := range resp.Locs {

		for _, s := range w.Locs {
			for _, r := range s.Locs {
				for _, b := range r.Locs {
					e.bins = append(e.bins, b.Uid)
				}
			}
		}
	}

}
