package api

import (
	"black-friday/fail"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// these typed errors map to gRPC status codes and automatic error handling logic
// extra context can be added via gRPC details to help human in debugging problems
// we verify status codes in tests, extra context is not verified
var (
	ErrNotUnimplemented = status.New(codes.Unimplemented, "Implement me!")
	ErrPrecondition     = status.New(codes.FailedPrecondition, "failed precondition")
	ErrArgument         = status.New(codes.InvalidArgument, "invalid argument")
	ErrBadMove          = status.New(codes.FailedPrecondition, "bad location move")
	ErrLocationNotFound = status.New(codes.NotFound, "location not found")
	ErrProductNotFound  = status.New(codes.NotFound, "product not found")

	ErrReservationNotFound = status.New(codes.NotFound, "reservation not found")
	ErrAlreadyExists       = status.New(codes.AlreadyExists, "already exists")
	ErrNotEnough           = status.New(codes.FailedPrecondition, "not enough quantity")
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
