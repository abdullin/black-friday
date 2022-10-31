package test

import (
	specs "black-friday/env/specs"
	"black-friday/inventory/api"
	"context"
	"fmt"
	"log"
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

	ctx := context.Background()
	env := specs.NewEnv(ctx, ":memory:")
	defer env.Close()

	env.EnsureSchema()

	// speed test

	oks, fails := 0, 0

	for i, s := range api.Specs {

		fmt.Printf("#%d. %s - taking too much time...", i+1, yellow(s.Name))

		tx, err := env.BeginTx()
		if err != nil {
			log.Panicln("begin tx", err)
		}

		result := env.RunSpec(s, tx)
		if err := tx.Rollback(); err != nil {
			log.Panicln("roll back")
		}

		deltas := result.Deltas

		fmt.Print(ERASE, "\r")
		if len(deltas) == 0 && err == nil {
			//fmt.Printf(" ✔\n")
			oks += 1
		} else {
			fails += 1
			specs.PrintFull(s, result)
			println()
		}

	}

	fmt.Printf("Total: ✔%d X%d\n", oks, fails)

}

func guard(err error) {
	if err != nil {
		panic(err)
	}
}
