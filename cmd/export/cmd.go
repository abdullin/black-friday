package export

import (
	"black-friday/env/specs"
	"black-friday/inventory/api"
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

type cmd struct {
}

func (c cmd) Help() string {
	return "Export specs into the text format"
}

func (c cmd) Run(args []string) int {

	var file string
	flags := flag.NewFlagSet("export", flag.ExitOnError)
	flags.StringVar(&file, "file", "export.txt", "File to export to.")

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}

	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("Can't create file %s: %s\n", file, err)
		return 1
	}
	defer f.Close()

	api.Sort()

	err = specs.WriteSpecs(api.Specs, f)
	if err != nil {
		fmt.Printf("Failed to write specs: %s\n", err)
		return 1
	}

	fmt.Printf("Wrote %d specs to %s\n", len(api.Specs), file)
	return 0
}

func (c cmd) Synopsis() string {
	return "Export specs outside"
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
