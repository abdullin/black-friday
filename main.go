package main

import (
	"black-friday/inventory/api"
	"black-friday/inventory/app"
	"black-friday/inventory/db"
	"black-friday/inventory/features/locations"
	"black-friday/inventory/features/products"
	"black-friday/inventory/features/stock"
	"black-friday/specs"
	"context"
	"database/sql"
	"fmt"
	"github.com/abdullin/go-seq"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"reflect"
	"runtime"

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

	var services []*app.App
	var wg sync.WaitGroup
	for i := 0; i < cores; i++ {
		dbs, err := sql.Open("sqlite3", file)
		guard(err)
		defer dbs.Close()

		guard(db.CreateSchema(dbs))

		svc := app.New(dbs)
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
				for _, s := range specs.Specs {
					local_count += 1
					result, err := RunSpec(ctx, svc, s)
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

	fmt.Printf("Discovered %d specs\n", len(specs.Specs))

	file := "/tmp/tests.sqlite"
	file = ":memory:"
	_ = os.Remove(file)

	dbs, err := sql.Open("sqlite3", file)
	guard(err)
	defer dbs.Close()

	guard(db.CreateSchema(dbs))

	svc := app.New(dbs)

	ctx := context.Background()

	// speed test

	for _, s := range specs.Specs {

		result, err := RunSpec(ctx, svc, s)
		deltas := result.Deltas
		if len(deltas) == 0 && err == nil {
			fmt.Printf("✔ %s️\n", s.Name)
		} else {
			fmt.Printf(red("x %s\n"), s.Name)

			s.Print()

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

func guard(err error) {
	if err != nil {
		panic(err)
	}
}

type SpecResult struct {
	EventCount int
	Deltas     seq.Issues
}

func dispatch(ctx *app.Context, m proto.Message) (r proto.Message, err error) {

	switch t := m.(type) {
	case *api.AddLocationsReq:
		r, err = locations.Add(ctx, t)
	case *api.AddProductsReq:
		r, err = products.Add(ctx, t)
	case *api.UpdateInventoryReq:
		r, err = stock.Update(ctx, t)
	case *api.ListLocationsReq:
		r, err = locations.List(ctx, t)
	case *api.GetLocInventoryReq:
		r, err = stock.Query(ctx, t)
	default:
		return nil, fmt.Errorf("missing dispatch for %v", reflect.TypeOf(m))
	}

	if r != nil && reflect.ValueOf(r).IsNil() {
		r = nil
	}
	return r, err
}

func RunSpec(ctx context.Context, a *app.App, spec *specs.S) (*SpecResult, error) {

	c, err := a.Begin(ctx)
	if err != nil {
		log.Panicln(err)
	}

	defer c.Rollback()

	for i, e := range spec.Given {
		err, fail := c.Apply(e)

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

	c.TestClear()

	actualResp, err := dispatch(c, spec.When)
	actualStatus, _ := status.FromError(err)
	var actualEvents []proto.Message
	if err == nil {
		actualEvents = c.TestGet()
	}

	eventCount += len(actualEvents)

	issues := spec.Compare(actualResp, actualStatus, actualEvents)

	return &SpecResult{
		EventCount: eventCount,
		Deltas:     issues,
	}, nil
}
