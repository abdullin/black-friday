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
	ErrNotUnimplemented = status.Error(codes.Unimplemented, "Implement me!")
	ErrPrecondition     = status.Error(codes.FailedPrecondition, "failed precondition")
	ErrArgument         = status.Errorf(codes.InvalidArgument, "invalid argument")

	ErrBadMove          = status.Error(codes.FailedPrecondition, "bad location move")
	ErrLocationNotFound = status.Error(codes.NotFound, "location not found")
	ErrProductNotFound  = status.Error(codes.NotFound, "product not found")
	ErrAlreadyExists    = status.Error(codes.AlreadyExists, "already exists")
	ErrNotEnough        = status.Error(codes.FailedPrecondition, "not enough quantity")
)

func ErrInternal(err error, code fail.Code) error {
	return status.Error(codes.Internal, fmt.Sprintf("fail-%d: %s", code, err))
}

func ErrArgNil(field string) error {
	return status.Error(codes.InvalidArgument, fmt.Sprintf("'%s' is nil", field))
}

func ErrArgInvalid(what string, value any, reason string) error {
	return status.Errorf(codes.InvalidArgument,
		"%v is not valid for %s: %s", value, what, reason)
}

func ErrInvalidOp(why string) error {
	return status.Error(codes.FailedPrecondition, why)

}
