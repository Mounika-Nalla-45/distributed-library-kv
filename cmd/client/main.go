package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "distributed-library-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewKVServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// PUT
	putResponse, err := client.Put(ctx, &pb.PutRequest{
		Key:   "B101",
		Value: "Computer Networks",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PUT success:", putResponse.GetSuccess())

	// GET
	getResponse, err := client.Get(ctx, &pb.GetRequest{
		Key: "B101",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GET found:", getResponse.GetFound())
	fmt.Println("GET value:", getResponse.GetValue())

	// DELETE
	deleteResponse, err := client.Delete(ctx, &pb.DeleteRequest{
		Key: "B101",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DELETE success:", deleteResponse.GetSuccess())

	// GET after DELETE
	getResponse, err = client.Get(ctx, &pb.GetRequest{
		Key: "B101",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GET after DELETE found:", getResponse.GetFound())
}
