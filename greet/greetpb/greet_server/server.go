package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc_tutorial/greet/greetpb"
	"log"
	"net"
)

type Server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *Server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet Function was called %v", req)
	firstname := req.GetGreeting().GetFirstName()
	result := "Hello " + firstname
	res := &greetpb.GreetResponse{Result: result}
	return res, nil
}

func main() {
	fmt.Println("Hello world")

	// create listener on default port 50051
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	// create grpc server
	s := grpc.NewServer()

	// Register a service
	greetpb.RegisterGreetServiceServer(s, &Server{})

	// bing port to grpc server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
