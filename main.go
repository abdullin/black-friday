package main

import (
	"context"
	"database/sql"
	"fmt"
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

	ctx := context.Background()

	for _, s := range tests.Specs {

		deltas, err := run_spec(ctx, svc, s)
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

func run_spec(ctx context.Context, svc *inventory.Service, spec *tests.Spec) ([]*seq.Delta, error) {

	tx := svc.GetTx(ctx)

	defer tx.Rollback()

	for _, e := range spec.Given {
		tx.Apply(e)
	}

	nested := context.WithValue(ctx, "tx", tx)
	actual, err := svc.Dispatch(nested, spec.When)

	actualStatus, _ := status.FromError(err)
	issues := seq.Diff(spec.ThenResponse, actual, "response")

	if spec.ThenError != actualStatus.Code() {
		actualErr := fmt.Sprintf("%s: %s", actualStatus.Code(), err.Error())

		issues = append(issues, &seq.Delta{
			Expected: spec.ThenError,
			Actual:   actualErr,
			Path:     "status",
		})
	}

	return issues, nil
}
