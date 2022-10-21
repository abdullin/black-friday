package stat

import (
	"black-friday/fail"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	DuplicateName = status.Error(codes.AlreadyExists, "Duplicate name")
	NotFound      = status.Error(codes.NotFound, "Entity not found")
)

func Internal(err error, code fail.Code) error {
	return status.Error(codes.Internal, fmt.Sprintf("fail-%d: %s", code, err))
}

func ArgNil(field string) error {
	return status.Error(codes.InvalidArgument, fmt.Sprintf("'%s' is nil", field))
}
