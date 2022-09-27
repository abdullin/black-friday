package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sdk-go/inventory"
	"sdk-go/protos"
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
	protos.RegisterInventoryServiceServer(grpcServer, inventory.NewService())

	// this allows to call this server with commands like:
	// grpcurl -plaintext localhost:8080 schema.InventoryService/ListLocation
	// grpcurl -plaintext localhost:8080 list
	reflection.Register(grpcServer)
	log.Panicln(grpcServer.Serve(lis))
}
