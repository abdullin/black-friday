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

			printSpec(s)

			if err != nil {
				fmt.Printf(red("  FATAL: %s\n"), err.Error())
			}

			for _, d := range deltas {
				fmt.Printf("  Δ %s\n", d.String())
			}
		}

	}

}

func printSpec(s *tests.Spec) {
	println(s.Name)
	if len(s.Given) > 0 {
		println("GIVEN:")
		for i, e := range s.Given {
			println(fmt.Sprintf("%d. %s", i+1, seq.Format(e)))
		}
	}
	println(fmt.Sprintf("WHEN: %s", seq.Format(s.When)))
	if s.ThenResponse != nil {
		println(fmt.Sprintf("THEN RESPONSE: %s", seq.Format(s.ThenResponse)))
	}
	if len(s.ThenEvents) > 0 {
		println("THEN EVENTS:")
		for i, e := range s.ThenEvents {
			println(fmt.Sprintf("%d. %s", i+1, seq.Format(e)))
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
	tx.TestClearEvents()

	nested := context.WithValue(ctx, "tx", tx)
	actualResp, err := svc.Dispatch(nested, spec.When)
	actualStatus, _ := status.FromError(err)
	actualEvents := tx.TestGetEvents()
	issues := seq.Diff(spec.ThenResponse, actualResp, "response")

	if len(actualEvents) != len(spec.ThenEvents) {
		issues = append(issues, &seq.Delta{
			Expected: spec.ThenEvents,
			Actual:   actualEvents,
			Path:     "events",
		})
	} else {
		for i, e := range spec.ThenEvents {
			p := fmt.Sprintf("events[%d]", i)
			issues = append(issues, seq.Diff(e, actualEvents[i], p)...)
		}
	}

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
