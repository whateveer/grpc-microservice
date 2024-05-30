package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/whateveer/payment-grpc/cmd/grpc-microservice/payment"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *server) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	fmt.Println("Received payment request")
	// Simulate payment processing
	transactionId := "1234567890"
	return &pb.PaymentResponse{
		TransactionId: transactionId,
		Message:       "Paid",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, &server{})

	log.Println("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
