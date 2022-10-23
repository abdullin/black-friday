package main

import (
	"black-friday/inventory"
	"context"
	"database/sql"
	"fmt"
	"github.com/abdullin/go-seq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"os"
	"reflect"
	"runtime"

	"black-friday/tests"
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
	if env := os.Getenv("REPL_ID"); env != "" {
		// REPLIT is limited by default
		cores = 1
	}

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
	var eventCount int64
	seconds := 1

	for i := 0; i < cores; i++ {
		go func(pos int) {
			svc := services[pos]
			duration := time.Second * time.Duration(seconds)
			var local_count int64
			var localEventCount int64
			for time.Since(started) < duration {
				for _, s := range tests.Specs {
					local_count += 1
					result, err := run_spec(ctx, svc, s)
					if err != nil {
						panic(err)
					}
					localEventCount += int64(result.EventCount)
				}
			}

			atomic.AddInt64(&count, local_count)
			atomic.AddInt64(&eventCount, localEventCount)
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Printf("running specs at %.1f kHz\n", float64(count)/1000.0/float64(seconds))
	fmt.Printf("applying events at %.1f kHz\n", float64(eventCount)/1000.0/float64(seconds))

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

		result, err := run_spec(ctx, svc, s)
		deltas := result.Deltas
		if len(deltas) == 0 && err == nil {
			fmt.Printf("✔ %s️\n", s.Name)
		} else {
			fmt.Printf(red("x %s\n"), s.Name)

			printSpec(s)

			if err != nil {
				fmt.Printf(red("  FATAL: %s\n"), err.Error())
			}

			fmt.Println(yellow("ISSUES:"))

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
			println(fmt.Sprintf("  %d. %s", i+1, seq.Format(e)))
		}
	}
	println(fmt.Sprintf("%s\n  %s", yellow("WHEN:"), seq.Format(s.When)))
	if s.ThenResponse != nil {
		println(fmt.Sprintf("%s\n  %s", yellow("THEN RESPONSE:"), seq.Format(s.ThenResponse)))
	}
	if s.ThenError != codes.OK {
		println(fmt.Sprintf("%s\n  %s", yellow("THEN ERROR:"), s.ThenError))
	}
	if len(s.ThenEvents) > 0 {
		println(yellow("THEN EVENTS:"))
		for i, e := range s.ThenEvents {
			println(fmt.Sprintf("  %d. %s", i+1, seq.Format(e)))
		}
	}
}

func guard(err error) {
	if err != nil {
		panic(err)
	}
}

type SpecResult struct {
	EventCount int
	Deltas     seq.Issues
}

func run_spec(ctx context.Context, svc *inventory.Service, spec *tests.Spec) (*SpecResult, error) {

	tx := svc.GetTx(ctx)

	defer tx.Rollback()

	for i, e := range spec.Given {
		err, fail := tx.Apply(e)
		if err != nil {
			panic(fmt.Sprintf("#%v problem with spec '%s' event %d.%s: %s",
				fail,
				spec.Name,
				i+1,
				reflect.TypeOf(e).String(),
				err))
		}
	}

	eventCount := len(spec.Given)

	tx.TestClearEvents()

	nested := context.WithValue(ctx, "tx", tx)
	actualResp, err := svc.Dispatch(nested, spec.When)
	actualStatus, _ := status.FromError(err)
	var actualEvents []proto.Message
	if err == nil {
		actualEvents = tx.TestGetEvents()
	}

	eventCount += len(actualEvents)

	issues := spec.Compare(actualResp, actualStatus, actualEvents)

	return &SpecResult{
		EventCount: eventCount,
		Deltas:     issues,
	}, nil
}
