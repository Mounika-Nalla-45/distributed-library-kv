package main

import (
	"distributed-library-kv/internal/cluster"
	"distributed-library-kv/internal/grpcserver"
	"distributed-library-kv/internal/raft"
	"distributed-library-kv/internal/storage"
	pb "distributed-library-kv/proto"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	nodeID := flag.String("id", "node1", "node ID")
	port := flag.String("port", "50051", "server port")
	dataDir := flag.String("data", "./data/node1", "data directory")

	flag.Parse()

	store, err := storage.NewLSMTree(*dataDir, 100)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	nodes := []cluster.Node{
		{
			ID:      "node1",
			Address: "localhost:50051",
		},
		{
			ID:      "node2",
			Address: "localhost:50052",
		},
		{
			ID:      "node3",
			Address: "localhost:50053",
		},
	}

	replicator := cluster.NewReplicator(nodes, *nodeID)
	raftNode := raft.NewNode(*nodeID)

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	service := grpcserver.NewServer(store, replicator, raftNode)
	pb.RegisterKVServiceServer(grpcServer, service)

	fmt.Printf(
		"Node %s running on port %s\n",
		*nodeID,
		*port,
	)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
