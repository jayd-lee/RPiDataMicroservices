package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/jayd-lee/RPiDataMicroservices/proto/sensor"
	"github.com/jayd-lee/RPiDataMicroservices/sensorservice/handler"
	"google.golang.org/grpc"
)

const (
	port        = ":50051"
	rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	maxRetries  = 5
	retryDelay  = 2 * time.Second
)

func main() {
	// Create TCP listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Try to create handler with retries
	var sensorHandler *handler.SensorHandler
	for i := 0; i < maxRetries; i++ {
		sensorHandler, err = handler.NewSensorHandler(rabbitmqURL)
		if err == nil {
			break
		}
		log.Printf("Failed to create handler (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}
	if err != nil {
		log.Fatalf("Failed to create handler after %d attempts: %v", maxRetries, err)
	}
	defer sensorHandler.Close()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	grpcServer := grpc.NewServer()
	pb.RegisterSensorServiceServer(grpcServer, sensorHandler)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down gracefully...")
	grpcServer.GracefulStop()
}
