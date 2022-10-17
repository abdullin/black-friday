package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"runtime"
	"sdk-go/inventory"
	"sdk-go/seq"
	"sdk-go/tests"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	CLEAR  = "\033[0m"
	RED    = "\033[91m"
	YELLOW = "\033[93m"

	GREEN = "\033[32m"

	ANOTHER = "\033[34m"
)

func red(s string) string {
	return fmt.Sprintf("%s%s%s", RED, s, CLEAR)
}
func yellow(s string) string {

	return fmt.Sprintf("%s%s%s", YELLOW, s, CLEAR)
}

func speed_test() {

	file := ":memory:"

	ctx := context.Background()

	cores := runtime.NumCPU()
	fmt.Printf("Speed test with %d cores... ", cores)

	var services []*inventory.Service
	var wg sync.WaitGroup
	for i := 0; i < cores; i++ {
		db, err := sql.Open("sqlite3", file)
		guard(err)
		defer db.Close()

		guard(inventory.CreateSchema(db))

		svc := inventory.NewService(db)
		services = append(services, svc)
		wg.Add(1)
	}

	// speed test

	started := time.Now()

	var count int64
	seconds := 1

	for i := 0; i < cores; i++ {
		go func(pos int) {
			svc := services[pos]
			duration := time.Second * time.Duration(seconds)
			var local_count int64
			for time.Since(started) < duration {
				for _, s := range tests.Specs {
					local_count += 1
					run_spec(ctx, svc, s)
				}
			}

			atomic.AddInt64(&count, local_count)
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Printf("running specs at %.1f kHz\n", float64(count)/1000.0/float64(seconds))

}

func main() {

	speed_test()

	fmt.Printf("Discovered %d specs\n", len(tests.Specs))

	file := "/tmp/tests.sqlite"
	file = ":memory:"
	_ = os.Remove(file)

	db, err := sql.Open("sqlite3", file)
	guard(err)
	defer db.Close()

	guard(inventory.CreateSchema(db))

	svc := inventory.NewService(db)

	ctx := context.Background()

	// speed test

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
				fmt.Printf("  %sΔ %s%s\n", ANOTHER, d.String(), CLEAR)
			}
			println()
		}

	}

}

func printSpec(s *tests.Spec) {
	//println(s.Name)
	if len(s.Given) > 0 {
		println(yellow("GIVEN:"))
		for i, e := range s.Given {
			println(fmt.Sprintf("%d. %s", i+1, seq.Format(e)))
		}
	}
	println(fmt.Sprintf("%s %s", yellow("WHEN:"), seq.Format(s.When)))
	if s.ThenResponse != nil {
		println(fmt.Sprintf("%s %s", yellow("THEN RESPONSE:"), seq.Format(s.ThenResponse)))
	}
	if s.ThenError != codes.OK {
		println(fmt.Sprintf("%s %s", yellow("THEN ERROR:"), s.ThenError))
	}
	if len(s.ThenEvents) > 0 {
		println(yellow("THEN EVENTS:"))
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

	for i, e := range spec.Given {
		err := tx.Apply(e)
		if err != nil {
			panic(fmt.Sprintf("Problem with spec '%s' precondition %d: %s", spec.Name, i+1, err))
		}
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
		actualErr := "OK"
		if actualStatus.Code() != codes.OK {

			actualErr = fmt.Sprintf("%s: %s", actualStatus.Code(), err)
		}

		issues = append(issues, &seq.Delta{
			Expected: spec.ThenError,
			Actual:   actualErr,
			Path:     "status",
		})
	}

	return issues, nil
}
