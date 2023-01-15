package subject

import (
	"black-friday/env/specs"
	"black-friday/inventory/api"
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"reflect"
)

type subject struct {
	env *specs.Env
	api.UnimplementedSpecServiceServer
}

func (s *subject) About(ctx context.Context, r *api.AboutRequest) (*api.AboutResponse, error) {
	return &api.AboutResponse{
		Author:  "Rinat Abdullin",
		Detail:  "golang",
		Contact: "@abdullin on twitter and mastodon.social",
	}, nil
}

func (s *subject) Spec(ctx context.Context, request *api.SpecRequest) (*api.SpecResponse, error) {

	tx, err := s.env.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	for i, e := range request.Given {
		var event proto.Message

		event, err := e.UnmarshalNew()
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal given[%d] %s: %w", i, e.GetTypeUrl(), err)
		}
		err, fail := tx.Apply(event, false)
		if err != nil {
			return nil, fmt.Errorf("#%v problem with given[%d] %d.%s: %w",
				fail,
				i,
				reflect.TypeOf(e).String(),
				err)
		}

		tx.ApplyModelEvent(e)
	}

	actualReq, err := request.When.UnmarshalNew()
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal when: %w", err)
	}

	tx.Events = nil
	resp, st := specs.Dispatch(tx, actualReq)

	if err := tx.Rollback(); err != nil {
		return nil, fmt.Errorf("fail to rollback: %w", err)
	}
	specResponse := &api.SpecResponse{}

	if st != nil {
		specResponse.Status = st.Proto().Code
		specResponse.Error = st.Message()

		// on error tx events will be ignored
		tx.Events = nil
	}
	if resp != nil {
		packed, err := anypb.New(resp)
		if err != nil {
			return nil, fmt.Errorf("fail to marshal response: %w", err)
		}
		specResponse.Response = packed
	}

	for _, e := range tx.Events {
		packed, err := anypb.New(e)
		if err != nil {
			return nil, fmt.Errorf("failed to pack event: %w", err)
		}

		specResponse.Events = append(specResponse.Events, packed)
	}

	return specResponse, nil
}
