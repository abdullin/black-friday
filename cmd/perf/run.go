package perf

import (
	specs2 "black-friday/env/specs"
	"black-friday/inventory/api"
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func speed_test(cores int, specs []*api.Spec, seconds int) {

	file := ":memory:"

	duration := time.Second * time.Duration(seconds)
	// set timeout, just in case
	ctx, cancel := context.WithTimeout(context.Background(), duration+time.Second)

	fmt.Printf("Speed test with %d core(s) for %ds\n", cores, seconds)

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
		specs, given int64
		dispatchTime int64
		givenTime    int64
	}
	var global counter

	for i := 0; i < cores; i++ {
		go func(pos int) {
			svc := services[pos]
			var local counter
			for time.Since(started) < duration {
				for _, s := range specs {
					local.specs += 1
					tx, err := svc.BeginTx()
					if err != nil {
						panic(err)
					}
					result := svc.RunSpec(s, tx)
					if err := tx.Rollback(); err != nil {
						panic(err)
					}
					local.given += int64(len(s.Given))
					local.givenTime += int64(result.GivenTime)

					local.dispatchTime += int64(result.Dispatch)
				}
			}

			atomic.AddInt64(&global.given, local.given)
			atomic.AddInt64(&global.specs, local.specs)
			atomic.AddInt64(&global.dispatchTime, local.dispatchTime)
			atomic.AddInt64(&global.givenTime, local.givenTime)
			wg.Done()
		}(i)
	}
	wg.Wait()
	cancel()

	hz := func(name string, count int64, op time.Duration) []string {
		khz := float64(count) / 1000.0 / op.Seconds()
		ops := int(float64(count) / op.Seconds())
		dur := op / time.Duration(count)

		return []string{
			name,
			fmt.Sprintf("%d", count),
			fmt.Sprintf("%d", ops),
			fmt.Sprintf("%.1f", khz),
			dur.String(),
		}

	}

	data := [][]string{
		hz("run spec", global.specs, duration),
		hz("apply event", global.given, time.Duration(global.givenTime)),
		hz("request", global.specs, time.Duration(global.dispatchTime)),
	}

	fmt.Println()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Operation", "Total", "ops/sec", "kHz", "sec per op"})

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render() // Send output

	fmt.Printf("\nexecuted %d specs\n\n", global.specs)
	fmt.Println("CAVEAT: This benchmarks event-driven spec tests with a minimal database. Real world performance will be worse due to: disk flush, DB growth, event store and network commit latency.")
}
