package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"p-bitbucket.imovetv.com/heracles/cmw/greet/greetpb"
)

func main() {
	fmt.Println("Hello, I am a client!")

	// Create connection
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	// When done, close it
	defer cc.Close()

	// Create client
	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Created client: %f\n", c)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Starting Unary RPC...")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Briggs",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}
