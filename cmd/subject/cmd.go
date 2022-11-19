package subject

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
	var addr string

	flags := flag.NewFlagSet("test", flag.ExitOnError)
	flags.StringVar(&db, "db", ":memory:", "sqlite db to use")
	flags.StringVar(&addr, "addr", "127.0.0.1:50051", "Subject to test")

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}

	serve_specs(db, addr)

	return 0
}

func (c cmd) Synopsis() string {
	return "run specs"
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
