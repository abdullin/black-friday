package export

import (
	"black-friday/env/specs"
	"black-friday/inventory/api"
	"github.com/mitchellh/cli"
)

type cmd struct {
}

func (c cmd) Help() string {
	return "Export specs into the text format"
}

func (c cmd) Run(args []string) int {

	for _, s := range api.Specs {
		println(specs.SpecToParseableString(s))
		println("=====================================")
	}

	return 0
}

func (c cmd) Synopsis() string {
	//TODO implement me
	panic("implement me")
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
