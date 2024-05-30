package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/whateveer/payment-grpc/cmd/grpc-microservice/payment"
	"google.golang.org/grpc"

	"github.com/jung-kurt/gofpdf"
	"gopkg.in/gomail.v2"
)

type server struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *server) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	fmt.Println("Received payment request")

	// Simulate payment processing
	transactionId := "1234567890"
	message := "Paid"

	// Generate PDF receipt
	receiptPath, err := generateReceipt(transactionId, req.FullName, req.Email)
	if err != nil {
		return nil, err
	}

	// Send receipt via email
	err = sendReceiptByEmail(req.Email, receiptPath)
	if err != nil {
		return nil, err
	}

	return &pb.PaymentResponse{
		TransactionId: transactionId,
		Message:       message,
	}, nil
}

func generateReceipt(transactionID, fullName, email string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "AITUNDER")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Transaction ID: %s", transactionID))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", time.Now().Format(time.RFC1123)))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Service: Premium account")
	pdf.Ln(10)
	pdf.Cell(40, 10, "Unit Price: $5")
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Customer: %s", fullName))
	receiptPath := fmt.Sprintf("%s_receipt.pdf", transactionID)
	err := pdf.OutputFileAndClose(receiptPath)
	if err != nil {
		return "", err
	}
	return receiptPath, nil
}

func sendReceiptByEmail(email, receiptPath string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "aitunderapp.notifications@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your Payment Receipt")
	m.SetBody("text/html", "Thank you for your payment. Please find your receipt attached.")
	m.Attach(receiptPath)

	d := gomail.NewDialer("smtp.gmail.com", 587, "aitunderapp.notifications@gmail.com", "your-email-password")

	return d.DialAndSend(m)
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
