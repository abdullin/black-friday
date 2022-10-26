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
)

func ErrInternal(err error, code fail.Code) error {
	return status.Error(codes.Internal, fmt.Sprintf("fail-%d: %s", code, err))
}

func ErrArgNil(field string) error {
	return status.Error(codes.InvalidArgument, fmt.Sprintf("'%s' is nil", field))
}
