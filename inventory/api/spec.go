package api

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

type Spec struct {
	Seq          int
	Name         string
	Given        []proto.Message
	When         proto.Message
	ThenResponse proto.Message
	ThenError    codes.Code
	ThenEvents   []proto.Message
}

var Specs []*Spec

func Define(s *Spec) {
	s.Seq = len(Specs) + 1
	Specs = append(Specs, s)
}
