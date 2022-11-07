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

	fmt.Println("Start simulation")

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

	fmt.Println("DURATION    DB SIZE  LOCATIONS   PRODUCTS     ON-HAND    RESERVED")
	for i := 0; i < 20; i++ {
		started := time.Now()

		if _, err := e.AddWarehouse(ctx); err != nil {
			log.Panicln(err)
		}

		e.AddProducts(ctx, 1000)
		e.AddInventory(ctx, 100)
		e.ReserveInventory(ctx, 1000)

		funcName(ctx, file, started, e, a)

	}
	return 0
}

func funcName(ctx context.Context, file string, started time.Time, e *env, a *node.Env) {

	tx, err := a.Begin(ctx)
	if err != nil {
		log.Panicln(err)
	}
	var onHand, reserved int64
	tx.QueryRow("SELECT SUM(OnHand) FROM Inventory")(&onHand)
	tx.QueryRow("SELECT SUM(Quantity) FROM Reserves")(&reserved)

	defer tx.Rollback()

	bytes := Size(file, file+"-wal")
	fmt.Printf("%5d ms %10s %10d %10d  %10d  %10d\n",
		time.Since(started).Milliseconds(), ByteCountDecimal(bytes), e.locations, e.products, onHand, reserved)
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
