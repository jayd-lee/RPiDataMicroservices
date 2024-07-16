package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	pb "github.com/jayd-lee/RPiDataMicroservices/proto/sensor"
	"github.com/streadway/amqp"
)

type SensorHandler struct {
	pb.UnimplementedSensorServiceServer
	rabbitmqConn *amqp.Connection
	rabbitmqChan *amqp.Channel
}

type SensorData struct {
	DeviceID    string    `json:"device_id"`
	Location    string    `json:"location"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewSensorHandler(rabbitmqURL string) (*SensorHandler, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare the queue
	_, err = ch.QueueDeclare(
		"sensor_data", // queue name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return nil, err
	}

	return &SensorHandler{
		rabbitmqConn: conn,
		rabbitmqChan: ch,
	}, nil
}

func (h *SensorHandler) StreamSensorData(ctx context.Context, data *pb.SensorData) (*pb.SensorResponse, error) {
	log.Printf("Received sensor data from device %s: Temperature: %.2f, Humidity: %.2f\n",
		data.DeviceId, data.Temperature, data.Humidity)

	// Prepare data for RabbitMQ
	sensorData := SensorData{
		DeviceID:    data.DeviceId,
		Location:    data.Location,
		Temperature: data.Temperature,
		Humidity:    data.Humidity,
		Timestamp:   time.Now(),
	}

	// Convert to JSON
	jsonData, err := json.Marshal(sensorData)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return nil, err
	}

	// Publish to RabbitMQ
	err = h.rabbitmqChan.Publish(
		"",            // exchange
		"sensor_data", // queue name
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		log.Printf("Error publishing to RabbitMQ: %v", err)
		return nil, err
	}

	return &pb.SensorResponse{
		Success: true,
		Message: "Data received and published successfully",
	}, nil
}

func (h *SensorHandler) Close() {
	if h.rabbitmqChan != nil {
		h.rabbitmqChan.Close()
	}
	if h.rabbitmqConn != nil {
		h.rabbitmqConn.Close()
	}
}
