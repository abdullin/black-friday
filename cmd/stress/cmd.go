package stress

import (
	"black-friday/env/node"
	"black-friday/env/pipe"
	"black-friday/inventory"
	"black-friday/inventory/api"
	"context"
	"fmt"
	"github.com/mitchellh/cli"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

type cmd struct{}

func (c cmd) Help() string {
	return "Run stress tests"
}

func (c cmd) Synopsis() string {
	return `Run a stress test operation against`
}

func (c cmd) Run(args []string) int {

	fmt.Println("Open simulation")

	ctx, stop := node.Cancel()
	defer stop()

	file := "/tmp/stress.sqlite"

	_ = os.Remove(file)

	a := node.NewEnv(ctx, file)
	a.EnsureSchema()

	// create server
	s := grpc.NewServer()

	server := inventory.New(a)
	api.RegisterInventoryServiceServer(s, server)

	fmt.Println("Setup simulated network")
	channel, cancel := pipe.ConnectToServer(ctx, s)
	defer cancel()

	client := api.NewInventoryServiceClient(channel)

	e := NewEnv(client)

	global := time.Now()
	for i := 0; i < 90; i++ {

		if i%20 == 0 {
			fmt.Println("#     TIME    DUR     DB      LOCs   SKUs  ON-HAND  RESERVE    SALES    REJECT   PENDING FULFILLED ENTITIES   EVENTS")
		}

		started := time.Now()

		e.AddWarehouse(ctx)

		e.AddProducts(ctx)

		for k := 0; k < 100; k++ {

			e.AddInventory(ctx)

			for s := 0; s < i*2+1; s++ {
				e.TrySell(ctx)
			}
		}

		e.TryFulfull(ctx, e.reservations.Len()*3/4)

		funcName(ctx, file, i, global, started, e, a)

		a.Bank.SaveSample(fmt.Sprintf("/tmp/trace_%02d.jsonl", i))
		a.Bank.SaveReport(fmt.Sprintf("/tmp/gross_%02d.txt", i))
		a.Bank.Clear()

		if time.Since(global) > time.Second*60 {
			//break
		}

	}
	return 0
}

func funcName(ctx context.Context, file string, i int, global, started time.Time, e *env, a *node.Env) {

	tx, err := a.Begin(ctx)
	if err != nil {
		log.Panicln(err)
	}
	var onHand, reserved, reservations, entities int64
	tx.QueryRow("SELECT SUM(OnHand) FROM Inventory")(&onHand)
	tx.QueryRow("SELECT ifnull(SUM(Quantity),0) FROM Reserves")(&reserved)
	tx.QueryRow("SELECT COUNT(*) FROM Reservations")(&reservations)
	tx.QueryRow("SELECT SUM(seq) FROM sqlite_sequence")(&entities)

	defer tx.Rollback()

	bytes := Size(file, file+"-wal")
	fmt.Printf("%2d %5dms %5dms %8s %6d %6d %8d %8d %8d  %8d  %8d %8d  %8d %8d\n",
		i,
		time.Since(global).Milliseconds(),
		time.Since(started).Milliseconds(),
		ByteCountDecimal(bytes),
		e.locations,
		e.products,
		onHand,
		reserved,
		e.sales,
		e.reject,
		reservations,
		e.fulfilled,
		entities,
		node.EventCount,
	)
}

func Size(names ...string) int64 {
	var size int64

	for _, x := range names {
		if s, err := os.Stat(x); err == nil {
			size += s.Size()
		}
	}
	return size
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
