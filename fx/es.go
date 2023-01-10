package fx

import (
	"context"
	"google.golang.org/protobuf/proto"
)

type EventStore interface {
	Publish(ctx context.Context, es []proto.Message) error
	Close() error
}
