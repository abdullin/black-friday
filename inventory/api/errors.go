package api

import (
	"black-friday/fail"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrDuplicateName    = status.Error(codes.AlreadyExists, "Duplicate name")
	ErrNotFound         = status.Error(codes.NotFound, "Entity not found")
	ErrNotUnimplemented = status.Error(codes.Unimplemented, "Implement me!")
	ErrPrecondition     = status.Error(codes.FailedPrecondition, "failed precondition")
	ErrArgument         = status.Errorf(codes.InvalidArgument, "invalid argument")
)

func ErrInternal(err error, code fail.Code) error {
	return status.Error(codes.Internal, fmt.Sprintf("fail-%d: %s", code, err))
}

func ErrSkuNotFound(sku string) error {
	return status.Error(codes.NotFound, fmt.Sprintf("sku %s not found", sku))
}

func ErrArgNil(field string) error {
	return status.Error(codes.InvalidArgument, fmt.Sprintf("'%s' is nil", field))
}
