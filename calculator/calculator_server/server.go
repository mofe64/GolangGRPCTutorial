package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc_tutorial/calculator/calculatorpb"
	"io"
	"log"
	"net"
)

type Server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (s *Server) CalculateSum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Calculate Sum function was called %v\n", req)
	num1 := req.GetNum1()
	num2 := req.GetNum2()
	sum := num1 + num2
	res := &calculatorpb.SumResponse{
		Sum: sum,
	}
	return res, nil

}

func (s *Server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("recieved prime number decomposition server stram rpc %v\n", req)
	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			err := stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			if err != nil {
				log.Fatalf("Error sending stream %v \n", err)
			}
			number = number / divisor
		} else {
			divisor++
		}
	}
	return nil
}

func (*Server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	sum := int32(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error reading stream %v\n", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func main() {
	fmt.Println("Hello World from calc server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}

}
