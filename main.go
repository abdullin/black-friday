package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"sdk-go/inventory"
	"sdk-go/seq"
	"sdk-go/tests"

	_ "github.com/mattn/go-sqlite3"
)

const (
	CLEAR = "\033[0m"
	RED   = "\033[91m"
)

func red(s string) string {
	return fmt.Sprintf("%s%s%s", RED, s, CLEAR)
}

func main() {

	fmt.Printf("Run %d specs\n", len(tests.Specs))

	file := "/tmp/tests.sqlite"
	file = ":memory:"
	_ = os.Remove(file)

	db, err := sql.Open("sqlite3", file)
	guard(err)
	defer db.Close()

	guard(inventory.CreateSchema(db))

	svc := inventory.NewService(db)

	for _, s := range tests.Specs {

		deltas, err := run_spec(svc, s)
		if len(deltas) == 0 && err == nil {
			fmt.Printf("✔ %s️\n", s.Name)
		} else {
			fmt.Printf(red("x %s\n"), s.Name)

			if err != nil {
				fmt.Printf(red("  FATAL: %s\n"), err.Error())
			}

			for _, d := range deltas {
				fmt.Printf("  %s\n", d.String())
			}
		}

	}

}

func guard(err error) {
	if err != nil {
		panic(err)
	}
}

func run_spec(svc *inventory.Service, spec *tests.Spec) ([]*seq.Delta, error) {

	ctx := context.Background()
	tx := svc.GetTx(ctx)

	defer tx.Rollback()

	for _, e := range spec.Given {
		tx.Apply(e)
	}

	actual, err := svc.Dispatch(context.WithValue(ctx, "tx", tx), spec.When)

	deltas1 := seq.Diff(spec.ThenResponse, actual, "response")

	actualStatus, _ := status.FromError(err)

	if spec.ThenError != actualStatus.Code() {

		var actualErr any
		if actualStatus.Code() == codes.Unknown {
			actualErr = fmt.Sprintf("Uknown:%s", err.Error())
		} else {
			actualErr = actualStatus.Code()
		}

		deltas1 = append(deltas1, &seq.Delta{
			Expected: spec.ThenError,
			Actual:   actualErr,
			Path:     "status",
		})
	}

	return deltas1, nil
}
