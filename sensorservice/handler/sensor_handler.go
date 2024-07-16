package handler

import (
	"context"
	"log"

	pb "github.com/jayd-lee/RPiDataMicroservices/proto/sensor"
	"github.com/streadway/amqp"
)

type SensorHandler struct {
	pb.UnimplementedSensorServiceServer
	rabbitmqConn *amqp.Connection
}

func NewSensorHandler(rabbitmqURL string) (*SensorHandler, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}

	return &SensorHandler{
		rabbitmqConn: conn,
	}, nil
}

func (h *SensorHandler) StreamSensorData(ctx context.Context, data *pb.SensorData) (*pb.SensorResponse, error) {
	log.Printf("Received sensor data from device %s: Temperature: %.2f, Humidity: %.2f\n",
		data.DeviceId, data.Temperature, data.Humidity)

	return &pb.SensorResponse{
		Success: true,
		Message: "Data received successfully",
	}, nil
}
