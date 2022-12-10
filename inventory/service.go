package inventory

import (
	"black-friday/fx"
	. "black-friday/inventory/api"
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
	UnimplementedInventoryServiceServer
}

func New(a fx.Transactor) InventoryServiceServer {

	return &server{
		app: a,
	}
}

func apiDispatch[A proto.Message, B proto.Message](s *server, c context.Context, req A, inner func(c fx.Tx, a A) (b B, st *status.Status)) (B, error) {
	var nilB B

	a := s.app
	ctx, err := a.Begin(c)

	trace := ctx.Trace()
	if !trace.Disabled() {
		trace.Begin(string(req.ProtoReflect().Descriptor().Name()))
	}

	if err != nil {
		return nilB, err
	}
	defer func() {
		err := ctx.Rollback()

		trace.End()
		trace.Close()

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

func (s *server) AddLocations(ctx context.Context, req *AddLocationsReq) (*AddLocationsResp, error) {
	return apiDispatch(s, ctx, req, locations.Add)
}
func (s *server) MoveLocation(ctx context.Context, req *MoveLocationReq) (*MoveLocationResp, error) {
	return apiDispatch(s, ctx, req, locations.Move)
}
func (s *server) ListLocations(ctx context.Context, req *ListLocationsReq) (*ListLocationsResp, error) {
	return apiDispatch(s, ctx, req, locations.List)
}

func (s *server) AddProducts(ctx context.Context, req *AddProductsReq) (*AddProductsResp, error) {
	return apiDispatch(s, ctx, req, products.Add)
}

func (s *server) UpdateInventory(c context.Context, r *UpdateInventoryReq) (*UpdateInventoryResp, error) {
	return apiDispatch(s, c, r, stock.Update)
}

func (s *server) GetLocInventory(c context.Context, r *GetLocInventoryReq) (*GetLocInventoryResp, error) {
	return apiDispatch(s, c, r, stock.Query)
}

func (s *server) Reserve(c context.Context, r *ReserveReq) (*ReserveResp, error) {
	return apiDispatch(s, c, r, stock.Reserve)
}

func (s *server) Fulfill(c context.Context, r *FulfillReq) (*FulfillResp, error) {
	return apiDispatch(s, c, r, stock.Fulfill)
}
