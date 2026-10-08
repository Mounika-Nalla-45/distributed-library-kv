package main

import (
	"distributed-library-kv/internal/cluster"
	"distributed-library-kv/internal/grpcserver"
	"distributed-library-kv/internal/raft"
	"distributed-library-kv/internal/storage"
	pb "distributed-library-kv/proto"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return fallback
}

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
			Address: "node1:50051",
		},
		{
			ID:      "node2",
			Address: "node2:50052",
		},
		{
			ID:      "node3",
			Address: "node3:50053",
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
	otherNodes := []string{}

	for _, node := range nodes {
		if node.ID != *nodeID {
			otherNodes = append(otherNodes, node.Address)
		}
	}

	go func() {
		time.Sleep(2 * time.Second)

		if *nodeID == "node1" {
			// Initial leader for the 3-node cluster.
			raftNode.BecomeLeader()

			fmt.Printf("Node %s became LEADER\n", *nodeID)

			// Leader sends heartbeats to node2 and node3.
			raftNode.StartHeartbeat(otherNodes)

			return
		}

		// node2 and node3 start as followers.
		fmt.Printf("Node %s is FOLLOWER\n", *nodeID)

		// Followers monitor the leader.
		go func() {
			for {
				time.Sleep(2 * time.Second)

				if raftNode.IsLeader() {
					continue
				}

				// Only start an election when heartbeat has timed out.
				if raftNode.ShouldStartElection() {
					fmt.Printf("Node %s starting election...\n", *nodeID)

					if raftNode.StartElection(otherNodes) {
						fmt.Printf("Node %s became LEADER\n", *nodeID)

						raftNode.StartHeartbeat(otherNodes)
						return
					}

					fmt.Printf(
						"Node %s election failed, retrying...\n",
						*nodeID,
					)
				}
			}
		}()
	}()
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
