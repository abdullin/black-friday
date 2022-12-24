package api

import (
	"fmt"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"strings"
)

type Spec struct {
	// from one to 5, how hard is this?
	Level        int
	Seq          int
	Name         string
	Comments     string
	Given        []proto.Message
	When         proto.Message
	ThenResponse proto.Message
	ThenError    *status.Status
	ThenEvents   []proto.Message
}

var Specs []*Spec

func Sort() {

	source := Specs
	var target []*Spec

	order := []proto.Message{

		&AddProductsReq{},
		&AddLocationsReq{},
		&ListLocationsReq{},
		&MoveLocationReq{},
		&GetLocInventoryReq{},
		&UpdateInventoryReq{},
		&ReserveReq{},
	}

	for _, o := range order {

		match := o.ProtoReflect().Descriptor().Name()
		for i, s := range source {
			if s == nil {
				continue
			}
			if s.When.ProtoReflect().Descriptor().Name() == match {
				target = append(target, s)
				source[i] = nil
			}
		}
	}

	for _, s := range source {
		if s != nil {
			target = append(target, s)
		}
	}
	Specs = target
}

func (s *Spec) ToTestString() string {

	var b strings.Builder
	ln := func(text string, args ...any) {
		_, err := b.WriteString(fmt.Sprintf(text, args...) + "\n")
		if err != nil {
			panic(err)
		}
	}

	ln(s.Name)

	if len(s.Given) > 0 {
		ln("GIVEN:")

		for _, e := range s.Given {
			ln("%s %s", e.ProtoReflect().Descriptor().Name(), prototext.Format(e))
		}
	}
	if s.When != nil {
		ln("WHEN: %s %s", s.When.ProtoReflect().Descriptor().Name(), prototext.Format(s.When))
	}
	if s.ThenResponse != nil {
		ln("THEN: %s %s", s.ThenResponse.ProtoReflect().Descriptor().Name(), prototext.Format(s.ThenResponse))
	}

	if len(s.ThenEvents) > 0 {
		ln("EVENTS:")
		for _, e := range s.ThenEvents {
			ln("%s %s", e.ProtoReflect().Descriptor().Name(), prototext.Format(e))
		}
	}
	if s.ThenError != nil {
		ln("ERROR: %s", s.ThenError)
	}

	return b.String()

}

func Define(s *Spec) {
	// cleanup the comments
	s.Comments = strings.Trim(s.Comments, "\n")

	s.Seq = len(Specs) + 1
	Specs = append(Specs, s)
}
