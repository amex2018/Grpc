package main

import (
	"context"
	"io"
	"log"
	"net"

	pb "github.com/VENOLD/grpc/grpc/product"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

// ProductServer is the server type for handling Product service requests
type ProductServer struct {
    pb.UnimplementedProductServer
    mongoClient *mongo.Client
}

// AddProducts is the implementation of the AddProducts RPC method
func (s *ProductServer) AddProducts(stream pb.Product_AddProductsServer) error {
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            // We've finished reading the stream. Send the final message.
            return stream.SendAndClose(&pb.ProductResponse{Result: "All products added successfully"})
        }
        if err != nil {
            return err
        }

        // Insert product data into MongoDB
        collection := s.mongoClient.Database("Product-grpc").Collection("products")
        _, err = collection.InsertOne(context.Background(), req)
        if err != nil {
            log.Printf("Error inserting into MongoDB: %v", err)
            return err
        }

        // Send a response after processing each request.
        if err := stream.Send(&pb.ProductResponse{Result: "Product added successfully"}); err != nil {
            return err
        }
    }
}

func main() {
    // Connect to MongoDB
    mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://127.0.0.1"))
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err := mongoClient.Disconnect(context.Background()); err != nil {
            log.Fatal(err)
        }
    }()

    // Create a gRPC server
    server := grpc.NewServer()

    // Register the ProductServer with the gRPC server
    pb.RegisterProductServer(server, &ProductServer{mongoClient: mongoClient})

    // Listen on a port
    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatal(err)
    }

    log.Println("gRPC server is running on :50051")
    // Serve the gRPC server
    if err := server.Serve(listener); err != nil {
        log.Fatal(err)
    }
}
