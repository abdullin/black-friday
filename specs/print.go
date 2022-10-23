package specs

import (
	"fmt"
	"github.com/abdullin/go-seq"
	"google.golang.org/grpc/codes"
)

func yellow(s string) string {

	return fmt.Sprintf("%s%s%s", YELLOW, s, CLEAR)
}

const (
	CLEAR  = "\033[0m"
	RED    = "\033[91m"
	YELLOW = "\033[93m"

	GREEN = "\033[32m"

	ANOTHER = "\033[34m"
)

func (s *S) Print() {
	//println(s.Name)
	if len(s.Given) > 0 {
		println(yellow("GIVEN:"))
		for i, e := range s.Given {
			println(fmt.Sprintf("  %d. %s", i+1, seq.Format(e)))
		}
	}
	println(fmt.Sprintf("%s\n  %s", yellow("WHEN:"), seq.Format(s.When)))
	if s.ThenResponse != nil {
		println(fmt.Sprintf("%s\n  %s", yellow("THEN RESPONSE:"), seq.Format(s.ThenResponse)))
	}
	if s.ThenError != codes.OK {
		println(fmt.Sprintf("%s\n  %s", yellow("THEN ERROR:"), s.ThenError))
	}
	if len(s.ThenEvents) > 0 {
		println(yellow("THEN EVENTS:"))
		for i, e := range s.ThenEvents {
			println(fmt.Sprintf("  %d. %s", i+1, seq.Format(e)))
		}
	}
}
