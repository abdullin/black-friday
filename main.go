package main

import (
	"black-friday/cmd/explore"
	"black-friday/cmd/export"
	"black-friday/cmd/perf"
	"black-friday/cmd/stress"
	"black-friday/cmd/subject"
	"black-friday/cmd/test"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/cli"
	"log"
	"os"
)

func main() {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"test":    test.Factory,
		"perf":    perf.Factory,
		"explore": explore.Factory,
		"stress":  stress.Factory,
		"export":  export.Factory,
		"subject": subject.Factory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
