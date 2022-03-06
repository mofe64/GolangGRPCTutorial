package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc_tutorial/calculator/calculatorpb"
	"io"
	"log"
	"math"
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

func (s *Server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		number := req.GetNumber()
		if number > max {
			max = number
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Number: max,
			})
			if sendErr != nil {
				return err
			}
		}
	}
}

func (s *Server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument, fmt.Sprintf("Received invalid number %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		Number: int32(math.Sqrt(float64(number))),
	}, nil

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
