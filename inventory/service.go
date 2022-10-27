package inventory

import (
	"black-friday/inventory/api"
	"black-friday/inventory/app"
	"black-friday/inventory/features/locations"
	"black-friday/inventory/features/products"
	"black-friday/inventory/features/stock"
	"context"
	"database/sql"
	"google.golang.org/protobuf/proto"
	"log"
)

// server implements GRPC server. It wires together all features
type server struct {
	app *app.App
	api.UnimplementedInventoryServiceServer
}

func New(db *sql.DB) api.InventoryServiceServer {

	a := app.New(db)

	return &server{app: a}
}

func apiDispatch[A proto.Message, B proto.Message](a *app.App, c context.Context, req A, x func(c *app.Context, a A) (b B, err error)) (B, error) {

	ctx, err := a.Begin(c)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := ctx.Rollback()
		if err != nil {
			log.Println("Additional error while rolling back: %s", err)
		}
	}()

	response, handleErr := x(ctx, req)
	if handleErr != nil {
		return nil, handleErr
	}
	commitErr := ctx.Commit()
	if commitErr != nil {
		return nil, commitErr
	}
	return response, nil

}

func (s *server) AddLocations(ctx context.Context, req *api.AddLocationsReq) (*api.AddLocationsResp, error) {
	return apiDispatch(s.app, ctx, req, locations.Add)
}
func (s *server) MoveLocation(ctx context.Context, req *api.MoveLocationReq) (*api.MoveLocationResp, error) {
	return apiDispatch(s.app, ctx, req, locations.Move)
}
func (s *server) ListLocations(ctx context.Context, req *api.ListLocationsReq) (*api.ListLocationsResp, error) {
	return apiDispatch(s.app, ctx, req, locations.List)
}

func (s *server) AddProducts(ctx context.Context, req *api.AddProductsReq) (*api.AddProductsResp, error) {
	return apiDispatch(s.app, ctx, req, products.Add)
}

func (s *server) UpdateInventory(c context.Context, r *api.UpdateInventoryReq) (*api.UpdateInventoryResp, error) {
	return apiDispatch(s.app, c, r, stock.Update)
}

func (s *server) GetLocInventory(c context.Context, r *api.GetLocInventoryReq) (*api.GetLocInventoryResp, error) {
	return apiDispatch(s.app, c, r, stock.Query)
}

func (s *server) Reserve(c context.Context, r *api.ReserveReq) (*api.ReserveResp, error) {
	return apiDispatch(s.app, c, r, stock.Reserve)
}
