package perf

import (
	"flag"
	"github.com/mitchellh/cli"
	"runtime"
)

type cmd struct {
}

func (c cmd) Help() string {
	return "Run performance tests"
	//TODO implement me
}

var (
	DEFAULT_CORES = runtime.NumCPU() / 2
)

func (c cmd) Run(args []string) int {

	var cores int

	flags := flag.NewFlagSet("perf", flag.ExitOnError)
	flags.IntVar(&cores, "cores", DEFAULT_CORES, "number of cores to use")

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}

	speed_test(cores)
	return 0
}

func (c cmd) Synopsis() string {
	return "runs speed test"
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
