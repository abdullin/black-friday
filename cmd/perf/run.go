package perf

import (
	specs2 "black-friday/env/specs"
	"black-friday/inventory/api"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func speed_test(cores int) {

	file := ":memory:"

	duration := time.Second
	// set timeout, just in case
	ctx, cancel := context.WithTimeout(context.Background(), duration+time.Second)

	fmt.Printf("Speed test with %d cores... \n", cores)

	var services []*specs2.Env
	var wg sync.WaitGroup
	for i := 0; i < cores; i++ {

		svc := specs2.NewEnv(ctx, file)
		defer svc.Close()

		svc.EnsureSchema()
		services = append(services, svc)
		wg.Add(1)
	}

	// speed test

	started := time.Now()

	type counter struct {
		specs, events int64
		dispatchTime  int64
	}
	var global counter

	for i := 0; i < cores; i++ {
		go func(pos int) {
			svc := services[pos]
			var local counter
			for time.Since(started) < duration {
				for _, s := range api.Specs {
					local.specs += 1
					tx, err := svc.BeginTx()
					if err != nil {
						panic(err)
					}
					result := svc.RunSpec(s, tx)
					if err := tx.Rollback(); err != nil {
						panic(err)
					}
					local.events += int64(result.EventCount)

					local.dispatchTime += int64(result.Dispatch)
				}
			}

			atomic.AddInt64(&global.events, local.events)
			atomic.AddInt64(&global.specs, local.specs)
			atomic.AddInt64(&global.dispatchTime, local.dispatchTime)
			wg.Done()
		}(i)
	}
	wg.Wait()
	cancel()

	hz := func(count int64, op time.Duration) string {
		khz := float64(count) / 1000.0 / op.Seconds()
		ops := int(float64(count) / op.Seconds())

		return fmt.Sprintf("%d ops/sec  (%.1f kHz)", ops, khz)
	}

	fmt.Printf("executed %d specs\n", global.specs)
	fmt.Printf("running specs:   %s\n", hz(global.specs, duration))
	fmt.Printf("applying events: %s\n", hz(global.events, duration))
	fmt.Printf("request speed:   %s\n", hz(global.specs, time.Duration(global.dispatchTime)))

}
