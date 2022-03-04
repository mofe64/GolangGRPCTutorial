package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc_tutorial/greet/greetpb"
	"log"
)

func main() {
	fmt.Println("Hello world from greet client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Could not close connection %v", err)
		}
	}(conn)
	c := greetpb.NewGreetServiceClient(conn)
	performUnaryOp(c)

}

func performUnaryOp(c greetpb.GreetServiceClient) {
	fmt.Println("Starting greet Unary RPC Op...")
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mofe",
			LastName:  "Ogunbiyi",
		},
	}

	greet, err := c.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("Error calling rpc: %v", err)
	}
	log.Printf("response from greet %v", greet.Result)
}
