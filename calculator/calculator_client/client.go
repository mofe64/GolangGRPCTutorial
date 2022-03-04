package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc_tutorial/calculator/calculatorpb"
	"log"
)

func main() {
	fmt.Println("Hello world from calc client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	),
	)
	if err != nil {
		log.Fatalf("Could not connecr %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Could not close connection %v", err)
		}
	}(conn)

	c := calculatorpb.NewCalculatorServiceClient(conn)
	performUnaryOp(c)
}

func performUnaryOp(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting calc Unary RPC Op...")
	request := &calculatorpb.SumRequest{
		Num1: 3,
		Num2: 10,
	}

	sum, err := c.CalculateSum(context.Background(), request)
	if err != nil {
		log.Fatalf("Error calling rpc %v", err)
	}
	log.Printf("response from calc %v", sum)
}
