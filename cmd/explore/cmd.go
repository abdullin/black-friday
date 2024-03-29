package explore

import (
	"black-friday/env/specs"
	"black-friday/env/uid"
	"black-friday/fx"
	"black-friday/inventory/api"
	"context"
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"log"
	"os"
)

type cmd struct {
}

func (c cmd) Help() string {
	return "preserve database state for exploration"
}

const DEFAULT_DB = "/tmp/debug.sqlite"

func (c cmd) Run(args []string) int {

	var specNum int
	var file string

	flags := flag.NewFlagSet("explore", flag.ExitOnError)
	flags.IntVar(&specNum, "spec", 1, "Spec id to explore")
	flags.StringVar(&file, "db", DEFAULT_DB, "Location to save db tp")

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}

	fx.Explore = true

	_ = os.Remove(file)
	ctx := context.Background()
	env := specs.NewEnv(file)

	defer env.Close()

	spec := api.Specs[specNum-1]
	uid.TestMode = true

	env.EnsureSchema()

	ttx, err := env.BeginTx(ctx)
	if err != nil {
		log.Panicln("begin tx:", err)
	}
	result := env.RunSpec(spec, ttx)
	specs.PrintFull(spec, result.Deltas)

	err = ttx.Commit()
	if err != nil {
		log.Panicln("tx: commit", err)
	}

	fmt.Println("Aggregate state saved to: ", file)

	return 0
}

func (c cmd) Synopsis() string {
	return "Explore "
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
