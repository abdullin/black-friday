package panda

import (
	"black-friday/fx"
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"strconv"
)

type store struct {
	client *kgo.Client
}

func (s *store) Publish(ctx context.Context, es []proto.Message) error {

	records := make([]*kgo.Record, 0, len(es))

	for _, e := range es {

		text := prototext.Format(e)
		header := string(e.ProtoReflect().Descriptor().Name())
		full := fmt.Sprintf("%s %s", header, text)

		record := kgo.KeyStringRecord("tenant-1", full)
		records = append(records, record)

	}
	//t.Trace().Arg(map[string]interface{}{"events": len(t.events), "bytes": byteCount})

	if err := s.client.BeginTransaction(); err != nil {
		panic(fmt.Errorf("error beginning transaction: %v\n", err))
	}

	results := s.client.ProduceSync(ctx, records...)
	if results.FirstErr() != nil {
		log.Panicln("Problem publishing", results.FirstErr())
	}

	// Attempt to commit the transaction and explicitly abort if the
	// commit was not attempted.
	switch err := s.client.EndTransaction(ctx, kgo.TryCommit); err {
	case nil:
	case kerr.OperationNotAttempted:
		panic("rollback")
	default:
		panic(fmt.Errorf("error committing transaction: %v\n", err))
	}
	return nil

}

func NewStore(topic string, seeds []string) (error, fx.EventStore) {

	//seeds := []string{"159.69.176.118:9092"}
	// One client can both produce and consume!
	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.

	producerId := strconv.FormatInt(int64(os.Getpid()), 10)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.TransactionalID(producerId),
		kgo.DefaultProduceTopic(topic),
	)
	if err != nil {
		panic(fmt.Errorf("Can't create kafka client %w", err))
	}

	return nil, &store{
		client: cl,
	}

}

func (es *store) Close() error {
	if es.client != nil {
		es.client.Close()
	}
	return nil
}
