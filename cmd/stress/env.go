package stress

import (
	"black-friday/env/rnd"
	"black-friday/env/uid"
	"black-friday/inventory/api"
	"context"
	"fmt"
	"log"
)

type env struct {
	products     int64
	locations    int64
	warehouses   int64
	reservations int64

	bins int64

	r *rnd.Rand

	client api.InventoryServiceClient
}

func NewEnv(client api.InventoryServiceClient) *env {
	return &env{client: client, r: rnd.New()}
}

func (e *env) ReserveInventory(ctx context.Context, count int) {

	for j := 0; j < count; j++ {

		e.reservations += 1
		name := fmt.Sprintf("sale-%d", e.reservations)

		c := int(e.r.Int63n(10) + 1)

		var items []*api.ReserveReq_Item

		for i := 0; i < c; i++ {
			product := e.r.Int63n(e.products-1) + 1
			items = append(items, &api.ReserveReq_Item{
				Sku:      SKU(product),
				Quantity: e.r.Int63n(5) + 1,
			})

		}

		_, _ = e.client.Reserve(ctx, &api.ReserveReq{
			Reservation: name,
			Items:       items,
		})
	}

}

func (e *env) AddInventory(ctx context.Context, count int) {

	for i := 0; i < count; i++ {

		product := e.r.Int63n(e.products-1) + 1
		locations := e.r.Int63n(e.locations-1) + 1

		quantity := e.r.Int63n(100)

		_, err := e.client.UpdateInventory(ctx, &api.UpdateInventoryReq{
			Location:     uid.Str(locations),
			Product:      uid.Str(product),
			OnHandChange: quantity,
		})

		if err != nil {
			log.Fatalln(err)
		}
	}

	// which product?

	//
}

func (e *env) AddProducts(ctx context.Context, count int) {
	var skus []string
	for p := 0; p < count; p++ {

		e.products += 1
		skus = append(skus, SKU(e.products))

	}

	prod := &api.AddProductsReq{Skus: skus}

	_, err := e.client.AddProducts(ctx, prod)
	if err != nil {
		log.Panicln(err)
	}
}

func SKU(e int64) string {
	return fmt.Sprintf("product-%d", e)
}

func (e *env) AddWarehouse(ctx context.Context) (*api.AddLocationsResp, error) {

	e.warehouses += 1

	whsName := fmt.Sprintf("WHS-%d", e.warehouses)

	whs := &api.AddLocationsReq_Loc{
		Name: whsName,
	}

	e.locations += 1

	// add rows
	for r := 0; r < 10; r++ {
		rowName := fmt.Sprintf("%s/ROW-%d", whsName, r+1)
		row := &api.AddLocationsReq_Loc{
			Name: rowName}
		whs.Locs = append(whs.Locs, row)

		e.locations += 1

		for s := 0; s < 5; s++ {
			shelfName := fmt.Sprintf("%s/SHELF-%d", rowName, s+1)
			shelf := &api.AddLocationsReq_Loc{Name: shelfName}
			row.Locs = append(row.Locs, shelf)

			e.locations += 1

			for b := 0; b < 9; b++ {
				e.bins += 1
				binName := fmt.Sprintf("BIN-%d", e.bins)
				bin := &api.AddLocationsReq_Loc{Name: binName}
				shelf.Locs = append(shelf.Locs, bin)

				e.locations += 1
			}

		}

	}

	return e.client.AddLocations(ctx, &api.AddLocationsReq{
		Locs:   []*api.AddLocationsReq_Loc{whs},
		Parent: uid.Str(0),
	})
}
