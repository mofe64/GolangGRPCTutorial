package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc_tutorial/greet/greetpb"
	"io"
	"log"
	"time"
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
	//performUnaryOp(c)
	//performServerStreamingOp(c)
	//performClientStreamingOp(c)
	performBiDirectionalStreamingOp(c)

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

func performServerStreamingOp(c greetpb.GreetServiceClient) {
	fmt.Println("Starting greet Server Streaming RPC Op...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mofe",
			LastName:  "Ogunbiyi",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling server stream rpc: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("End of stream ...")
			break
		}
		if err != nil {
			log.Fatalf("Error reading stream %v", err)
		}
		log.Printf("res from server streaming %v", msg.GetResult())

	}

}
func performClientStreamingOp(c greetpb.GreetServiceClient) {
	fmt.Println("Starting greet Client Streaming RPC Op...")
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mofe",
				LastName:  "Test",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mofe",
				LastName:  "Test",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Sam",
				LastName:  "Test",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tade",
				LastName:  "Test",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tayo",
				LastName:  "Test",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling long greet %v\n", err)
	}

	for _, req := range requests {
		fmt.Printf("sending req %v\n", req)
		err := stream.Send(req)
		time.Sleep(100 * time.Millisecond)
		if err != nil {
			log.Fatalf("Error sending req %v\n", err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receieving res %v\n", err)
	}
	fmt.Printf("res %v\n", res)
}

func performBiDirectionalStreamingOp(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Bi directional Streaming RPC Op...")

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mofe",
				LastName:  "Test",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mofe",
				LastName:  "Test",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Sam",
				LastName:  "Test",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tade",
				LastName:  "Test",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Tayo",
				LastName:  "Test",
			},
		},
	}
	// create stream
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream %v\n", err)
	}

	// Create a wait channel
	waitC := make(chan struct{})

	// send messages to server in go routine
	go func() {
		for _, req := range requests {
			fmt.Println("Attempting to send message")
			err := stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
			if err != nil {
				log.Fatalf("Error sending message %v\n", err)
				return
			}
			fmt.Println("Message sent")
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("Error closing stream %v\n", err)
			return
		}
	}()

	// receive messages from server in go routine
	go func() {
		// receive data from server
		for {
			res, err := stream.Recv()
			// if eof encountered i.e. No more data is being sent by the server
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving message %v\n", err)
			}
			fmt.Printf("received %v\n", res.GetResult())
		}
		// close our wait channel and unblock function
		close(waitC)
	}()

	//block until all operations are done
	<-waitC

}
