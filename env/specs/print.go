package specs

import (
	"black-friday/inventory/api"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"strings"
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

func Print(s *api.Spec) {
	//println(s.Name)
	if len(s.Given) > 0 {
		println(yellow("GIVEN:"))
		for i, e := range s.Given {
			println(fmt.Sprintf("  %d. %s", i+1, Format(e)))
		}
	}
	println(fmt.Sprintf("%s\n  %s", yellow("WHEN:"), Format(s.When)))
	if s.ThenResponse != nil {
		println(fmt.Sprintf("%s\n  %s", yellow("THEN RESPONSE:"), Format(s.ThenResponse)))
	}
	if s.ThenError != codes.OK {
		println(fmt.Sprintf("%s\n  %s", yellow("THEN ERROR:"), s.ThenError))
	}
	if len(s.ThenEvents) > 0 {
		println(yellow("THEN EVENTS:"))
		for i, e := range s.ThenEvents {
			println(fmt.Sprintf("  %d. %s", i+1, Format(e)))
		}
	}
}
func Format(val any) string {
	if val == nil {
		return "<nil>"
	}
	switch v := val.(type) {
	case proto.Message:

		repr := prototext.MarshalOptions{Multiline: false}.Format(v)
		return string(v.ProtoReflect().Descriptor().Name()) + " " + repr + ""
	case []proto.Message:
		names := []string{}
		for _, m := range v {
			names = append(names, string(m.ProtoReflect().Descriptor().Name()))
		}
		return fmt.Sprintf("[%s]", strings.Join(names, ", "))

	case error:
		return fmt.Sprintf("Error '%v'", v.Error())
	default:
		return fmt.Sprintf("'%v'", v)
	}
}
