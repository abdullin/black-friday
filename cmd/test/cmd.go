package test

import (
	"flag"
	"github.com/mitchellh/cli"
)

type cmd struct {
}

func (c cmd) Help() string {

	return "run specs"
}

func (c cmd) Run(args []string) int {

	var db string

	flags := flag.NewFlagSet("test", flag.ExitOnError)
	flags.StringVar(&db, "db", ":memory:", "sqlite db to use")

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}

	test_specs()
	return 0
}

func (c cmd) Synopsis() string {
	return "run specs"
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
