package inventory

import (
	"black-friday/fx"
	"black-friday/inventory/api"
	"black-friday/inventory/features/locations"
	"black-friday/inventory/features/products"
	"black-friday/inventory/features/stock"
	"context"
	"database/sql"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
)

// server implements GRPC server. It wires together all features
type server struct {
	app fx.Transactor
	api.UnimplementedInventoryServiceServer
}

func New(a fx.Transactor) api.InventoryServiceServer {

	return &server{app: a}
}

func apiDispatch[A proto.Message, B proto.Message](a fx.Transactor, c context.Context, req A, inner func(c fx.Tx, a A) (b B, st *status.Status)) (B, error) {
	var nilB B

	ctx, err := a.Begin(c)

	if err != nil {
		return nilB, err
	}
	defer func() {
		err := ctx.Rollback()
		if err == sql.ErrTxDone {
			return
		}
		if err != nil {
			log.Printf("Additional error while rolling back: %s\n", err)
		}
	}()

	response, st := inner(ctx, req)
	if st != nil {
		return nilB, st.Err()
	}

	commitErr := ctx.Commit()
	if commitErr != nil {
		return nilB, commitErr
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
