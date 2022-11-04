package perf

import (
	"black-friday/inventory/api"
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"log"
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
	var seconds int
	var exclude, include string

	flags := flag.NewFlagSet("perf", flag.ExitOnError)
	flags.IntVar(&cores, "cores", DEFAULT_CORES, "number of cores to use")

	flags.IntVar(&seconds, "sec", 1, "seconds to run the test")
	flags.StringVar(&exclude, "exclude", "", "spec exclude")
	flags.StringVar(&include, "include", "", "spec exclude")

	if err := flags.Parse(args); err != nil {
		flags.Usage()
		return 1
	}

	var specs []*api.Spec
	var includeRegex, excludeRegex *regexp.Regexp

	if len(include) > 0 {
		includeRegex = regexp.MustCompile(include)
	}
	if len(exclude) > 0 {
		excludeRegex = regexp.MustCompile(exclude)
	}

	for _, s := range api.Specs {
		includeThis := includeRegex == nil || includeRegex.MatchString(s.Name)
		excludeThis := excludeRegex != nil && excludeRegex.MatchString(s.Name)

		if includeThis && !excludeThis {
			specs = append(specs, s)
		}
	}
	fmt.Printf("Matched %d specs out of %d\n", len(specs), len(api.Specs))
	if len(specs) == 0 {
		log.Fatalln("No specs to work with!")
	}

	speed_test(cores, specs, seconds)
	return 0
}

func (c cmd) Synopsis() string {
	return "runs speed test"
}

func Factory() (cli.Command, error) {
	return &cmd{}, nil
}
