package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"net"
	"sdk-go/protos"
	"sdk-go/tests"
)

type Server struct {
	protos.UnimplementedSpecServiceServer
}

func (s *Server) Start(ctx context.Context, req *protos.StartReq) (*protos.StartResp, error) {
	log.Println("Starting session")

	return &protos.StartResp{
		Session: "",
		Task:    specToTask(0),
	}, nil
}

func (s *Server) Exec(ctx context.Context, req *protos.ExecReq) (*protos.ExecResp, error) {

	log.Println("Got result ", req.Result.Seq)

	spec := tests.Specs[req.Result.Seq]

	resp := mustUnwrap(req.Result.Response)
	status := status.New(codes.Code(req.Result.StatusCode), req.Result.Message)
	var events []proto.Message
	for _, e := range req.Result.Events {
		events = append(events, mustUnwrap(e))
	}
	issues := spec.Compare(resp, status, events)

	if len(issues) > 0 {
		f := &protos.Failure{
			Name: spec.Name,
		}
		for _, i := range issues {
			f.Issue = append(f.Issue, i.String())
		}
		return &protos.ExecResp{
			Results: &protos.Report{
				Success: int64(req.Result.Seq),
				Failed:  1,
				Fails:   []*protos.Failure{f},
			},
		}, nil
	}

	next := int(req.Result.Seq + 1)

	if next < len(tests.Specs)-1 {

		return &protos.ExecResp{
			Task:    specToTask(next),
			Results: nil,
		}, nil
	} else {
		return &protos.ExecResp{
			Results: &protos.Report{
				Success: int64(len(tests.Specs)),
				Failed:  0,
			},
		}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protos.RegisterSpecServiceServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
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

func specToTask(seq int) *protos.SpecTask {

	s := tests.Specs[seq]
	t := &protos.SpecTask{
		Seq:  uint32(seq),
		When: mustWrap(s.When),
	}
	for _, e := range s.Given {
		t.Given = append(t.Given, mustWrap(e))
	}
	return t

}
