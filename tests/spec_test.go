package tests

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/status"
	"sdk-go/inventory"
	. "sdk-go/protos"
	"sdk-go/seq"
	"testing"
)

func run_spec(t *testing.T, spec *Spec) {
	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	db, err := sql.Open("sqlite3", ":memory:")
	check(err)
	defer db.Close()

	check(inventory.CreateSchema(db))

	s := inventory.NewService(db)

	check(s.ApplyEvents(spec.Given))

	actual, err := s.Dispatch(context.Background(), spec.When)

	deltas1 := seq.Diff(spec.ThenResponse, actual, "response")

	actualStatus, _ := status.FromError(err)

	if spec.ThenError != actualStatus.Code() {
		deltas1 = append(deltas1, &seq.Delta{
			Expected: spec.ThenError,
			Actual:   actualStatus.Code(),
			Path:     "status",
		})
	}

	for _, d := range deltas1 {
		t.Error(d.String())
	}
}

func Test_Spec(t *testing.T) {

	for _, s := range Specs {
		t.Run(s.Name, func(t *testing.T) {
			run_spec(t, s)
		})
	}

}
