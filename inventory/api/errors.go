package api

import (
	"black-friday/fail"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// these are slighly more detailed error codes
// they are STILL very generic. TODO: add context details later
var (
	ErrNotUnimplemented = status.New(codes.Unimplemented, "Implement me!")
	ErrPrecondition     = status.New(codes.FailedPrecondition, "failed precondition")
	ErrArgument         = status.New(codes.InvalidArgument, "invalid argument")
	ErrBadMove          = status.New(codes.FailedPrecondition, "bad location move")
	ErrLocationNotFound = status.New(codes.NotFound, "location not found")
	ErrProductNotFound  = status.New(codes.NotFound, "product not found")
	ErrAlreadyExists    = status.New(codes.AlreadyExists, "already exists")
	ErrNotEnough        = status.New(codes.FailedPrecondition, "not enough quantity")
)

func ErrInternal(err error, code fail.Code) *status.Status {
	return status.New(codes.Internal, fmt.Sprintf("fail-%d: %s", code, err))
}

func ErrArgNil(field string) *status.Status {
	return status.New(codes.InvalidArgument, fmt.Sprintf("'%s' is nil", field))
}

func ErrArgInvalid(what string, value any, reason string) *status.Status {
	return status.Newf(codes.InvalidArgument,
		"%v is not valid for %s: %s", value, what, reason)
}

func ErrInvalidOp(why string) *status.Status {
	return status.New(codes.FailedPrecondition, why)

}
