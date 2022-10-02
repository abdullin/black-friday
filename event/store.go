package event

import (
	"context"
	"google.golang.org/protobuf/proto"
)

type Stream struct {
	Stream    string
	Signature uint64
}

type Result struct {
	Signature uint64
}

type Metadata struct {
}

type Reader interface {
	Close() error
	Recv() (proto.Message, Metadata, error)
}

type Store interface {
	Append(c context.Context, s *Stream, events ...proto.Message) (*Result, error)
	Load(c context.Context, s *Stream, count int) Reader
}
