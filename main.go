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
	RED   = "\033[91m"
)

func red(s string) string {
	return fmt.Sprintf("%s%s%s", RED, s, CLEAR)
}

func main() {

	fmt.Printf("Run %d specs\n", len(tests.Specs))

	for _, s := range tests.Specs {

		deltas, err := run_spec(s)
		if len(deltas) == 0 && err == nil {
			fmt.Printf("✔ %s️\n", s.Name)
		} else {
			fmt.Printf(red("x %s\n"), s.Name)

			if err != nil {
				fmt.Printf(red("  %s\n"), err.Error())
			}

			for _, d := range deltas {
				fmt.Printf("  %s\n", d.String())
			}
		}

	}

}

type SpecResult struct {
	Deltas []*seq.Delta
	Panic  error
}

func run_spec(spec *tests.Spec) ([]*seq.Delta, error) {
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

	err = s.ApplyEvents(spec.Given)

	if err != nil {
		return nil, err
	}

	actual, err := s.Dispatch(context.Background(), spec.When)
	if err != nil {
		return nil, err
	}

	deltas1 := seq.Diff(spec.ThenResponse, actual, "response")

	actualStatus, _ := status.FromError(err)

	if spec.ThenError != actualStatus.Code() {
		deltas1 = append(deltas1, &seq.Delta{
			Expected: spec.ThenError,
			Actual:   actualStatus.Code(),
			Path:     "status",
		})
	}

	return deltas1, nil
}
