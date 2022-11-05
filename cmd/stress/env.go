package stress

import (
	"black-friday/inventory/api"
	"context"
	"fmt"
	"log"
)

type env struct {
	products   int
	locations  int
	warehouses int

	client api.InventoryServiceClient
}

func (e *env) AddProducts(ctx context.Context, count int) {
	var skus []string
	for p := 0; p < count; p++ {

		e.products += 1
		skus = append(skus, sku(e.products))

	}

	prod := &api.AddProductsReq{Skus: skus}

	_, err := e.client.AddProducts(ctx, prod)
	if err != nil {
		log.Panicln(err)
	}
}

func (e *env) AddWarehouse(ctx context.Context) (*api.AddLocationsResp, error) {

	e.warehouses += 1

	name := fmt.Sprintf("WHS-%d", e.warehouses)

	whs := &api.AddLocationsReq_Loc{
		Name: name,
	}

	bins := 0

	// add rows
	for r := 0; r < 10; r++ {
		name := fmt.Sprintf("WHS-%d/ROW-%d", e.warehouses, r+1)
		row := &api.AddLocationsReq_Loc{
			Name: name}
		whs.Locs = append(whs.Locs, row)

		for s := 0; s < 5; s++ {
			shelfName := fmt.Sprintf("SHELF-%d", s+1)
			shelf := &api.AddLocationsReq_Loc{Name: shelfName}
			row.Locs = append(row.Locs, shelf)

			for b := 0; b < 9; b++ {
				bins += 1
				binName := fmt.Sprintf("BIN-%d", bins)
				bin := &api.AddLocationsReq_Loc{Name: binName}
				shelf.Locs = append(shelf.Locs, bin)
			}

		}

	}

	return e.client.AddLocations(ctx, &api.AddLocationsReq{
		Locs:   []*api.AddLocationsReq_Loc{whs},
		Parent: 0,
	})
}

func sku(id int) string {
	return fmt.Sprintf("product-%d", id)

}
