package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"p-bitbucket.imovetv.com/heracles/cmw/greet/greetpb"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet was called with %v\n", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello " + firstName + " " + lastName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello world!")
	// Create connection
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// Create gRPC server
	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	// Bind port to gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
