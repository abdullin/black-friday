package main

import (
	"flag"
	"fmt"
	"go-lite/inventory"
	"go-lite/schema"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	port := flag.Uint("port", 8080, "port")

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	schema.RegisterInventoryServiceServer(grpcServer, inventory.NewService())
	log.Panicln(grpcServer.Serve(lis))
}
