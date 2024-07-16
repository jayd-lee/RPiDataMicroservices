package main

/*
#cgo CFLAGS: -I${SRCDIR}/../csensor
#cgo LDFLAGS: ${SRCDIR}/../csensor/libsensor.a
#include "sensor.h"
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jayd-lee/RPiDataMicroservices/sensor/gosensor/device"
	"github.com/jayd-lee/RPiDataMicroservices/sensor/gosensor/publisher"
)

// Define publishers as package-level variables to be accessible in collectDataForDevice
var (
	grpcPub *publisher.GRPCPublisher
	rmqPub  *publisher.RabbitMQPublisher
)

func collectData() (float64, float64) {
	var temperature C.double
	var humidity C.double
	C.generate_sensor_data(&temperature, &humidity)
	return float64(temperature), float64(humidity)
}

func collectDataForDevice(device *device.Device, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			temp, hum := collectData()
			device.UpdateData(temp, hum)
			data := device.GetLastData()

			// Publish to gRPC
			if err := grpcPub.PublishData(device.ID, device.Location, temp, hum); err != nil {
				log.Printf("Failed to publish to gRPC: %v", err)
			}

			// Publish to RabbitMQ
			if err := rmqPub.PublishData(data); err != nil {
				log.Printf("Failed to publish to RabbitMQ: %v", err)
			}

			fmt.Printf("[Device: %s, Location: %s] Temperature: %.2f°C, Humidity: %.2f%%, Time: %s\n",
				device.ID,
				device.Location,
				data.Temperature,
				data.Humidity,
				data.Timestamp.Format(time.RFC3339))

		case <-stopCh:
			fmt.Printf("Stopping data collection for device %s\n", device.ID)
			return
		}
	}
}

func main() {
	// Initialize device manager
	deviceManager := device.NewDeviceManager()

	// Initialize publishers
	var err error
	grpcPub, err = publisher.NewGRPCPublisher("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create gRPC publisher: %v", err)
	}

	rmqPub, err = publisher.NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ publisher: %v", err)
	}
	defer rmqPub.Close()

	// Initialize devices
	devices := []struct {
		id       string
		location string
	}{
		{"device1", "Room 101"},
		{"device2", "Room 102"},
		{"device3", "Room 103"},
	}

	for _, d := range devices {
		if err := deviceManager.AddDevice(d.id, d.location); err != nil {
			log.Fatalf("Failed to add device: %v", err)
		}
	}

	// Create channels and WaitGroup for graceful shutdown
	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	// Set up signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Start collection for each device
	fmt.Println("Starting Go sensor data collection...")
	for _, device := range deviceManager.GetAllDevices() {
		wg.Add(1)
		go collectDataForDevice(device, stopCh, &wg)
	}

	// Wait for interrupt signal
	<-signalCh
	fmt.Println("\nReceived interrupt signal. Shutting down...")

	// Signal all goroutines to stop
	close(stopCh)

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Data collection stopped")
}
