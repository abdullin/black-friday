package perf

import "github.com/mitchellh/cli"

type cmd struct {
}

func (c cmd) Help() string {
	return "Run performance tests"
	//TODO implement me
}

func (c cmd) Run(args []string) int {
	speed_test()
	return 0
}

func (c cmd) Synopsis() string {
	return "runs speed test"
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
