package inventory

import (
	"errors"
	"github.com/mattn/go-sqlite3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func re[M proto.Message](m M, err error) (M, error) {

	if err == nil {
		return m, nil
	}

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		switch sqliteErr.Code {
		case sqlite3.ErrConstraint:
			return m, status.Error(codes.FailedPrecondition, "Constraint violation")
		default:
			return m, status.Errorf(codes.Internal, err.Error())
		}
	}

	return m, err
}
