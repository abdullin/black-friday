package tests

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"sdk-go/seq"
)

func (spec *Spec) Compare(resp proto.Message, status *status.Status, events []proto.Message) []*seq.Delta {

	issues := seq.Diff(spec.ThenResponse, resp, "response")
	if spec.ThenError != status.Code() {
		actualErr := "OK"
		if status.Code() != codes.OK {

			actualErr = fmt.Sprintf("%s: %s", status.Code(), status.Message())
		}

		issues = append(issues, &seq.Delta{
			Expected: spec.ThenError,
			Actual:   actualErr,
			Path:     "status",
		})

	}

	if len(events) != len(spec.ThenEvents) {
		issues = append(issues, &seq.Delta{
			Expected: spec.ThenEvents,
			Actual:   events,
			Path:     "events",
		})
	} else {
		for i, e := range spec.ThenEvents {
			p := fmt.Sprintf("events[%d]", i)
			issues = append(issues, seq.Diff(e, events[i], p)...)
		}
	}
	return issues
}
