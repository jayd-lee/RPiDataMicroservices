package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jayd-lee/RPiDataMicroservices/dbservice/handler"
)

const (
	dbConnStr   = "postgres://username:password@localhost:5432/sensordb?sslmode=disable"
	rabbitmqURL = "amqp://guest:guest@localhost:5672/"
)

func main() {
	dbHandler, err := handler.NewDBHandler(dbConnStr, rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to create DB handler: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := dbHandler.StartConsumingMessages(); err != nil {
			log.Fatalf("Failed to start consuming messages: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down...")
}
