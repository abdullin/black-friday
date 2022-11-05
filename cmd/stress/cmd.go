package stress

import (
	"black-friday/env/node"
	"black-friday/env/pipe"
	"black-friday/inventory"
	"black-friday/inventory/api"
	"fmt"
	"github.com/mitchellh/cli"
	"google.golang.org/grpc"
	"log"
	"os"
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

	fmt.Println("Client")
	client := api.NewInventoryServiceClient(channel)

	fmt.Println("List locations")
	r, e := client.AddProducts(ctx, &api.AddProductsReq{Skus: []string{"p1", "p2"}})
	if e != nil {
		log.Fatalln(e)
	}
	fmt.Println(r.String())
	return 0
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
