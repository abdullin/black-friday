package perf

import (
	specs2 "black-friday/env/specs"
	"black-friday/inventory/api"
	"black-friday/inventory/db"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

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

	var services []*specs2.Env
	var wg sync.WaitGroup
	for i := 0; i < cores; i++ {
		dbs, err := sql.Open("sqlite3", file)
		if err != nil {
			log.Panicln(err)
		}
		defer dbs.Close()

		err = db.CreateSchema(dbs)
		if err != nil {
			log.Panicln(err)
		}

		svc := specs2.NewEnv(ctx, dbs)
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
