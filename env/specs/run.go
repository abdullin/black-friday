package specs

import (
	"black-friday/inventory/api"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"
)

func (env *Env) RunSpec(spec *api.Spec) (*SpecResult, error) {

	ttx, err := env.Begin()
	if err != nil {
		return nil, fmt.Errorf("create tx: %w", err)
	}
	defer func() {
		err := ttx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			log.Panicln(fmt.Errorf("spec rollback: %w", err))
		}
	}()

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
	/*
		if spec.Name == "reservation reduced availability" {
			ttx.Commit()
			os.Exit(1)
		}

	*/
	actualStatus, _ := status.FromError(err)
	var actualEvents []proto.Message
	if err == nil {
		actualEvents = ttx.events
	}

	eventCount += len(actualEvents)

	issues := Compare(spec, actualResp, actualStatus, actualEvents)

	return &SpecResult{
		EventCount: eventCount,
		Deltas:     issues,
	}, nil
}
