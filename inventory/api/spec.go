package api

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

type Spec struct {
	Name         string
	Given        []proto.Message
	When         proto.Message
	ThenResponse proto.Message
	ThenError    codes.Code
	ThenEvents   []proto.Message
}

var Specs []*Spec

func Define(s *Spec) {
	Specs = append(Specs, s)
}
