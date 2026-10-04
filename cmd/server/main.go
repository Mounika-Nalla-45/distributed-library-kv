package main

import (
	"fmt"
	"log"
	"net"

	"distributed-library-kv/internal/grpcserver"
	"distributed-library-kv/internal/storage"
	pb "distributed-library-kv/proto"

	"google.golang.org/grpc"
)

func main() {
	// Create LSM storage.
	store, err := storage.NewLSMTree("./data", 100)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	// Create TCP listener.
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	// Create gRPC server.
	grpcServer := grpc.NewServer()

	// Register our KV service.
	service := grpcserver.NewServer(store)
	pb.RegisterKVServiceServer(grpcServer, service)

	fmt.Println("gRPC server running on port 50051")

	// Start server.
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
