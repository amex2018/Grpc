package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/VENOLD/grpc/grpc/product"
	"google.golang.org/grpc"
)

func main() {
    // Set up a connection to the server.
    conn, err := grpc.Dial(":50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }
    defer conn.Close()

    // Create a gRPC client
    client := pb.NewProductClient(conn)

    // Create a stream to the server
    stream, err := client.AddProducts(context.Background())
    if err != nil {
        log.Fatalf("Error calling AddProducts: %v", err)
    }

    // Send requests to the server
    requests := []*pb.ProductRequest{
        {ProductName: "SampleProduct1"},
        {ProductName: "SampleProduct2"},
        // Add more products as needed
    }

    for _, req := range requests {
        if err := stream.Send(req); err != nil {
            log.Fatalf("%v.Send(%v) = %v", stream, req, err)
        }
    }
    stream.CloseSend()

    // Receive responses from the server
    for {
        resp, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("%v.Recv() got error %v, want %v", stream, err, nil)
        }
        fmt.Printf("Received response: %s\n", resp.Result)
    }
}
