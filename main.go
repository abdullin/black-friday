package main

import (
	"black-friday/inventory/api"
	"black-friday/inventory/db"
	"black-friday/specs"
	"context"
	"database/sql"
	"fmt"
	"os"
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

	seconds := 1.0
	// set timeout, just in case
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(seconds+0.5))

	cores := runtime.NumCPU()
	if env := os.Getenv("REPL_ID"); env != "" {
		// REPLIT is limited by default
		cores = 1
	}

	fmt.Printf("Speed test with %d cores... ", cores)

	var services []*specs.Env
	var wg sync.WaitGroup
	for i := 0; i < cores; i++ {
		dbs, err := sql.Open("sqlite3", file)
		guard(err)
		defer dbs.Close()

		guard(db.CreateSchema(dbs))

		svc := specs.NewEnv(ctx, dbs)
		services = append(services, svc)
		wg.Add(1)
	}

	// speed test

	started := time.Now()

	var count int64
	var eventCount int64

	for i := 0; i < cores; i++ {
		go func(pos int) {
			svc := services[pos]
			duration := time.Second * time.Duration(seconds)
			var local_count int64
			var localEventCount int64
			for time.Since(started) < duration {
				for _, s := range api.Specs {
					local_count += 1
					result, err := svc.RunSpec(s)
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
	cancel()

	fmt.Printf("running specs at %.1f kHz\n", float64(count)/1000.0/float64(seconds))
	fmt.Printf("applying events at %.1f kHz\n", float64(eventCount)/1000.0/float64(seconds))

}

func main() {

	//speed_test()

	fmt.Printf("Discovered %d specs\n", len(api.Specs))

	file := "/tmp/tests.sqlite"
	file = ":memory:"
	_ = os.Remove(file)

	dbs, err := sql.Open("sqlite3", file)
	guard(err)
	defer dbs.Close()

	guard(db.CreateSchema(dbs))

	ctx := context.Background()
	env := specs.NewEnv(ctx, dbs)

	// speed test

	oks, fails := 0, 0

	for _, s := range api.Specs {

		fmt.Printf(s.Name)

		result, err := env.RunSpec(s)
		deltas := result.Deltas
		if len(deltas) == 0 && err == nil {
			fmt.Printf(" ✔\n")
			oks += 1
		} else {
			fails += 1
			fmt.Printf(red(" x\n"))

			specs.Print(s)

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

	fmt.Printf("Total: ✔%d x%d\n", oks, fails)

}

func guard(err error) {
	if err != nil {
		panic(err)
	}
}
