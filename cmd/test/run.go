package test

import (
	"black-friday/env/pipe"
	specs "black-friday/env/specs"
	"black-friday/env/uid"
	"black-friday/inventory/api"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
)

const (
	CLEAR  = "\033[0m"
	RED    = "\033[91m"
	YELLOW = "\033[93m"

	GREEN = "\033[32m"

	ANOTHER = "\033[34m"
	ERASE   = "\033[2K"
)

func red(s string) string {
	return fmt.Sprintf("%s%s%s", RED, s, CLEAR)
}
func yellow(s string) string {

	return fmt.Sprintf("%s%s%s", YELLOW, s, CLEAR)
}

func green(s string) string {

	return fmt.Sprintf("%s%s%s", GREEN, s, CLEAR)
}

func mustAny(p proto.Message) *anypb.Any {
	r, err := anypb.New(p)
	if err != nil {
		log.Panicln("failed to convert to any: %w", err)
	}
	return r
}

func mustMsg(a *anypb.Any) proto.Message {
	if a == nil {
		return nil
	}
	p, err := a.UnmarshalNew()
	if err != nil {
		log.Panicln("failed to convert from any: %w", err)
	}
	return p
}

func test_specs(db, addr string) {

	//speed_test()

	fmt.Printf("Found %d specs to run\n", len(api.Specs))

	var conn *grpc.ClientConn
	var err error

	// setup subject
	ctx := context.Background()
	if addr != "" {
		conn, err = grpc.DialContext(ctx, addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Panicln(err)
		}
	} else {
		env := specs.NewEnv(db)
		defer env.Close()

		env.EnsureSchema()

		subj := &subject{env: env}

		s := grpc.NewServer()
		api.RegisterSpecServiceServer(s, subj)

		fmt.Println("Setup simulated network")
		var cancel func()
		conn, cancel = pipe.ConnectToServer(ctx, s)
		defer cancel()
	}

	// setup client
	client := api.NewSpecServiceClient(conn)

	uid.TestMode = true

	// speed test

	oks, fails, issues := 0, 0, 0

	for i, s := range api.Specs {

		fmt.Printf("#%d. %s - taking too much time...", i+1, yellow(s.Name))

		request := &api.SpecRequest{
			When: mustAny(s.When),
		}

		for _, e := range s.Given {
			request.Given = append(request.Given, mustAny(e))
		}

		resp, err := client.Spec(ctx, request)

		if err != nil {
			fails += 1
			fmt.Printf("%s%s%s\n", RED, err.Error(), CLEAR)
		}
		var events []proto.Message
		for _, e := range resp.Events {
			events = append(events, mustMsg(e))
		}

		st := status.New(codes.Code(resp.Status), resp.Error)

		deltas := specs.Compare(s, mustMsg(resp.Response), st, events)
		issues += len(deltas)

		fmt.Print(ERASE, "\r")
		if len(deltas) == 0 && err == nil {
			//fmt.Printf(" ✔\n")
			oks += 1
		} else {
			fails += 1

			specs.PrintFull(s, deltas)
			println()
		}

	}

	fmt.Printf("PASS:%d FAIL:%d ISSUES:%d\n", oks, fails, issues)

}
