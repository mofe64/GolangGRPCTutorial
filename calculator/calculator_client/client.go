package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"grpc_tutorial/calculator/calculatorpb"
	"io"
	"log"
	"time"
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
	//performUnaryOp(c)
	//performServerStreamingOp(c)
	//performClientStreamingOp(c)
	//performBiDirectionalStreamingOp(c)
	performUnaryOp2(c)

}

func performUnaryOp2(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting calc Unary RPC Op...")
	request := &calculatorpb.SquareRootRequest{
		Number: -10,
	}
	sqrt, err := c.SquareRoot(context.Background(), request)
	if err != nil {
		respErr, ok := status.FromError(err)
		// ok signifies user error occurred
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("Invalid arg")
				return
			}
		} else {
			log.Fatalf("Internal Error calling rpc %v", err)
			return
		}

	}
	log.Printf("response from calc %v", sqrt.GetNumber())
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

func performServerStreamingOp(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting calculator Server Streaming RPC Op...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
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
		log.Printf("res from server streaming %v", msg.GetPrimeFactor())

	}
}

func performClientStreamingOp(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error opening stream %v\n", err)
	}
	numbers := []int32{3, 5, 9, 54, 23}
	for _, number := range numbers {
		err := stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
		if err != nil {
			log.Fatalf("Error sending stream %v\n", err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error closing stream %v\n", err)
	}
	fmt.Printf("Average is %v\n", res.Average)
}
func performBiDirectionalStreamingOp(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error opening stream %v\n", err)
	}
	waitC := make(chan struct{})
	go func() {
		numbers := []int32{4, 7, 6, 34, 12, 19, 45, 34, 12, 43, 65, 12, 5, 6, 8, 0, 123, 5, 7, 8, 678}
		for _, num := range numbers {
			err := stream.Send(&calculatorpb.FindMaximumRequest{
				Number: num,
			})
			time.Sleep(1000 * time.Millisecond)
			if err != nil {
				log.Fatalf("Error sending data  %v\n", err)
			}
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalf("Error closing stream %v\n", err)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data %v\n", err)
			}
			max := res.GetNumber()
			fmt.Printf("Received new max of .... %v\n", max)
		}
		close(waitC)
	}()

	<-waitC
}
