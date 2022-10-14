package tests

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/status"
	"sdk-go/inventory"
	"sdk-go/seq"
	"testing"
)

func guard(err error) {
	if err != nil {
		panic(err)
	}
}

func run_spec(t *testing.T, spec *Spec, s *inventory.Service) {

	tx := s.GetTx(context.Background())
	for _, e := range spec.Given {
		tx.Apply(e)
	}

	nested := context.WithValue(context.Background(), "tx", tx)

	actual, err := s.Dispatch(nested, spec.When)

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

	db := must(sql.Open("sqlite3", ":memory:"))
	defer db.Close()

	guard(inventory.CreateSchema(db))

	svc := inventory.NewService(db)

	for _, s := range Specs {
		t.Run(s.Name, func(t *testing.T) {
			run_spec(t, s, svc)
		})
	}

}
