package main

import (
	"context"
	"database/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"sdk-go/inventory"
	"sdk-go/protos"
)

func guard(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func mustWrap(msg proto.Message) *anypb.Any {
	if nil == msg {
		return nil
	}
	a, err := anypb.New(msg)
	if err != nil {
		log.Panicln("Failed to Marhsal")
	}
	return a
}
func mustUnwrap(a *anypb.Any) proto.Message {
	if nil == a {
		return nil
	}
	msg, err := a.UnmarshalNew()
	if err != nil {
		log.Panicln("Failed to unmarshal Any")
	}
	return msg
}

func runSpec(svc *inventory.Service, task *protos.SpecTask) *protos.SpecResult {

	tx := svc.GetTx(context.Background())

	defer tx.Rollback()

	for _, r := range task.Given {
		err, _ := tx.Apply(mustUnwrap(r))
		guard(err)

	}

	tx.TestClearEvents()
	nested := context.WithValue(context.Background(), "tx", tx)
	when := mustUnwrap(task.When)
	actualResp, err := svc.Dispatch(nested, when)

	r := &protos.SpecResult{
		Seq: task.Seq,
	}
	if err != nil {
		stat, _ := status.FromError(err)
		r.StatusCode = int32(stat.Code())
		r.Message = stat.Message()
	} else {
		for _, e := range tx.TestGetEvents() {
			r.Events = append(r.Events, mustWrap(e))
		}
		r.Response = mustWrap(actualResp)
	}
	return r
}

func main() {

	log.Println("Connecting to remote spec server")
	file := ":memory:"
	db, err := sql.Open("sqlite3", file)
	guard(err)
	defer db.Close()
	guard(inventory.CreateSchema(db))

	svc := inventory.NewService(db)
	ctx := context.Background()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := protos.NewSpecServiceClient(conn)

	start, err := client.Start(ctx, &protos.StartReq{})
	guard(err)
	log.Println(start.Session)

	task := start.Task
	guard(err)

	for {
		result := runSpec(svc, task)
		eval, err := client.Exec(ctx, &protos.ExecReq{
			Session: start.Session,
			Result:  result,
		})
		guard(err)
		if eval.Results != nil {
			log.Println(eval.Results)
			if eval.Results.Failed == 0 {
				log.Println("All specs passed!")
			}
			break
		}
		task = eval.Task
	}
}
