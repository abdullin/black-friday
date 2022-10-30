package test

import (
	specs2 "black-friday/env/specs"
	"black-friday/inventory/api"
	"black-friday/inventory/db"
	"context"
	"database/sql"
	"fmt"
	"github.com/abdullin/go-seq"
	"os"
	"strings"
)

const (
	CLEAR  = "\033[0m"
	RED    = "\033[91m"
	YELLOW = "\033[93m"

	GREEN = "\033[32m"

	ANOTHER = "\033[34m"
	ERASE   = "\033[2K"
)

func red(s string) string {
	return fmt.Sprintf("%s%s%s", RED, s, CLEAR)
}
func yellow(s string) string {

	return fmt.Sprintf("%s%s%s", YELLOW, s, CLEAR)
}

func green(s string) string {

	return fmt.Sprintf("%s%s%s", GREEN, s, CLEAR)
}

func test_specs() {

	//speed_test()

	fmt.Printf("Found %d specs to run\n", len(api.Specs))

	file := "/tmp/tests.sqlite"
	file = ":memory:"
	_ = os.Remove(file)

	dbs, err := sql.Open("sqlite3", file)
	guard(err)
	defer dbs.Close()

	guard(db.CreateSchema(dbs))

	ctx := context.Background()
	env := specs2.NewEnv(ctx, dbs)

	// speed test

	oks, fails := 0, 0

	for i, s := range api.Specs {

		fmt.Printf("#%d. %s - taking too much time...", i+1, yellow(s.Name))

		result, err := env.RunSpec(s)
		deltas := result.Deltas

		fmt.Print(ERASE, "\r")
		if len(deltas) == 0 && err == nil {
			//fmt.Printf(" ✔\n")
			oks += 1
		} else {
			fails += 1
			fmt.Printf(red("X %s\n"), red(s.Name))

			specs2.Print(s)

			if err != nil {
				fmt.Printf(red("  FATAL: %s\n"), err.Error())
			}

			fmt.Println(yellow("ISSUES:"))

			for _, d := range deltas {
				fmt.Printf("  %sΔ %s%s\n", ANOTHER, IssueToString(d), CLEAR)
			}
			println()
		}

	}

	fmt.Printf("Total: ✔%d X%d\n", oks, fails)

}

func IssueToString(d seq.Issue) string {
	return fmt.Sprintf("Expected %v to be %v but got %v",
		strings.Replace(seq.JoinPath(d.Path), ".[", "[", -1),
		specs2.Format(d.Expected),
		specs2.Format(d.Actual))

}

func guard(err error) {
	if err != nil {
		panic(err)
	}
}
