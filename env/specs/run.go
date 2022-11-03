package specs

import (
	"black-friday/inventory/api"
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func (env *Env) RunSpec(spec *api.Spec, ttx *tx) *SpecResult {

	for i, e := range spec.Given {
		err, fail := ttx.Apply(e)

		if err != nil {
			panic(fmt.Sprintf("#%v problem with spec '%s' event %d.%s: %s",
				fail,
				spec.Name,
				i+1,
				reflect.TypeOf(e).String(),
				err))
		}
	}

	eventCount := len(spec.Given)

	ttx.events = nil

	actualResp, err := dispatch(ttx, spec.When)
	var actualEvents []proto.Message
	if err == nil {
		actualEvents = ttx.events
	}

	eventCount += len(actualEvents)

	issues := Compare(spec, actualResp, err, actualEvents)

	return &SpecResult{
		EventCount: eventCount,
		Deltas:     issues,
	}
}
