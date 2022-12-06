package subject

import (
	specs "black-friday/env/specs"
	"black-friday/env/uid"
	"black-friday/inventory/api"
	"google.golang.org/grpc"
	"log"
	"net"
)

func serve_specs(db, addr string) {

	//speed_test()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	env := specs.NewEnv(db)
	defer env.Close()

	env.EnsureSchema()

	subj := &subject{env: env}

	s := grpc.NewServer()
	api.RegisterSpecServiceServer(s, subj)
	uid.TestMode = true

	log.Println("Serving on ", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
