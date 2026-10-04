package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "distributed-library-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := flag.String("address", "localhost:50051", "server address")
	operation := flag.String("operation", "put", "put, get or delete")
	key := flag.String("key", "B101", "key")
	value := flag.String("value", "Computer Networks", "value")

	flag.Parse()

	conn, err := grpc.NewClient(
		*address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewKVServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Connected to:", *address)

	switch *operation {

	case "put":
		response, err := client.Put(ctx, &pb.PutRequest{
			Key:   *key,
			Value: *value,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("PUT success:", response.GetSuccess())

	case "get":
		response, err := client.Get(ctx, &pb.GetRequest{
			Key: *key,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("GET found:", response.GetFound())
		fmt.Println("GET value:", response.GetValue())

	case "delete":
		response, err := client.Delete(ctx, &pb.DeleteRequest{
			Key: *key,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("DELETE success:", response.GetSuccess())

	default:
		log.Fatalf("unknown operation: %s", *operation)
	}
}
