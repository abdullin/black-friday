package perf

import (
	"black-friday/inventory/api"
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"regexp"
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
	var filter string

	flags := flag.NewFlagSet("perf", flag.ExitOnError)
	flags.IntVar(&cores, "cores", DEFAULT_CORES, "number of cores to use")
	flags.StringVar(&filter, "filter", "", "spec filter")

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}

	var regex *regexp.Regexp

	specs := api.Specs
	if len(filter) > 0 {

		specs = nil
		regex = regexp.MustCompile(filter)
		fmt.Printf("Filtering specs with %s\n", filter)

		for _, s := range api.Specs {
			if regex.MatchString(s.Name) {
				fmt.Printf("  - %s\n", s.Name)
				specs = append(specs, s)
			}
		}
	} else {

		fmt.Printf("Using %d specs in test\n", len(specs))
	}

	speed_test(cores, specs)
	return 0
}

func (c cmd) Synopsis() string {
	return "runs speed test"
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
