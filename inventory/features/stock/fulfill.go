package stock

import (
	"black-friday/fx"
	"black-friday/inventory/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Fulfill(ctx fx.Tx, req *api.FulfillReq) (*api.FulfillResp, *status.Status) {

	return nil, status.New(codes.Unimplemented, "TODO")

}
