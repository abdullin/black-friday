package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/status"
	"sdk-go/inventory"
	"sdk-go/seq"
	"sdk-go/tests"

	_ "github.com/mattn/go-sqlite3"
)

const (
	CLEAR = "\033[0m"
	RED   = "\033[91;4m"
)

func red(s string) string {
	return fmt.Sprintf("%s%s%s", RED, s, CLEAR)
}

func main() {

	fmt.Printf("Run %d specs\n", len(tests.Specs))

	for _, s := range tests.Specs {

		deltas := run_spec(s)
		if len(deltas) == 0 {
			fmt.Printf("✔ %s️\n", s.Name)
		} else {
			fmt.Printf(red("x %s\n"), s.Name)

			for _, d := range deltas {
				fmt.Printf("  %s\n", d.String())
			}
		}

	}

}

func run_spec(spec *tests.Spec) []*seq.Delta {
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

	return deltas1
}
