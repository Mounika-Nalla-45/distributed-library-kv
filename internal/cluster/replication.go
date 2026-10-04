package cluster

import (
	"context"
	"fmt"
	"time"

	pb "distributed-library-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Replicator struct {
	nodes []Node
	self  string
}

func NewReplicator(nodes []Node, self string) *Replicator {
	return &Replicator{
		nodes: nodes,
		self:  self,
	}
}

func (r *Replicator) ReplicatePut(key, value string) {
	for _, node := range r.nodes {
		if node.ID == r.self {
			continue
		}

		go r.sendPut(node, key, value)
	}
}

func (r *Replicator) ReplicateDelete(key string) {
	for _, node := range r.nodes {
		if node.ID == r.self {
			continue
		}

		go r.sendDelete(node, key)
	}
}

func (r *Replicator) sendPut(node Node, key, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		node.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer conn.Close()

	client := pb.NewKVServiceClient(conn)

	response, err := client.ReplicatePut(ctx, &pb.PutRequest{
		Key:   key,
		Value: value,
	})

	if err != nil {
		fmt.Println("Replication PUT failed:", node.ID, err)
		return
	}

	fmt.Println("Replication PUT successful:", node.ID, response.GetSuccess())
}

func (r *Replicator) sendDelete(node Node, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		node.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer conn.Close()

	client := pb.NewKVServiceClient(conn)

	_, _ = client.ReplicateDelete(ctx, &pb.DeleteRequest{
		Key: key,
	})
}
