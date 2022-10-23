package specs

import (
	"fmt"
	"github.com/abdullin/go-seq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type S struct {
	Name         string
	Given        []proto.Message
	When         proto.Message
	ThenResponse proto.Message
	ThenError    codes.Code
	ThenEvents   []proto.Message
}

func (spec *S) Compare(resp proto.Message, status *status.Status, events []proto.Message) seq.Issues {

	issues := seq.Diff(spec.ThenResponse, resp, "response")
	if spec.ThenError != status.Code() {
		actualErr := "OK"
		if status.Code() != codes.OK {

			actualErr = fmt.Sprintf("%s: %s", status.Code(), status.Message())
		}

		issues = append(issues, seq.Issue{
			Expected: spec.ThenError,
			Actual:   actualErr,
			Path:     []string{"status"},
		})

	}

	if len(events) != len(spec.ThenEvents) {
		issues = append(issues, seq.Issue{
			Expected: spec.ThenEvents,
			Actual:   events,
			Path:     []string{"events"},
		})
	} else {
		for i, e := range spec.ThenEvents {
			p := []string{"events", fmt.Sprintf("[%d]", i)}
			issues = append(issues, seq.Diff(e, events[i], p...)...)
		}
	}
	return issues
}
