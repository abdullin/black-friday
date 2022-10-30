package test

import "github.com/mitchellh/cli"

type cmd struct {
}

func (c cmd) Help() string {

	return "run specs"
}

func (c cmd) Run(args []string) int {
	test_specs()
	return 0
}

func (c cmd) Synopsis() string {
	return "run specs"
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
