package publisher

import (
	"context"
	"time"

	pb "github.com/jayd-lee/RPiDataMicroservices/proto/sensor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

type GRPCPublisher struct {
	client pb.SensorServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCPublisher(serverAddr string) (*GRPCPublisher, error) {
	// Configure connection with retry
	conn, err := grpc.Dial(serverAddr,
		grpc.WithInsecure(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				Multiplier: 1.5,
				MaxDelay:   2 * time.Second,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewSensorServiceClient(conn)
	return &GRPCPublisher{
		client: client,
		conn:   conn,
	}, nil
}

func (p *GRPCPublisher) PublishData(deviceID, location string, temp, humidity float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := p.client.StreamSensorData(ctx, &pb.SensorData{
		DeviceId:    deviceID,
		Location:    location,
		Temperature: temp,
		Humidity:    humidity,
		Timestamp:   time.Now().Format(time.RFC3339),
	})
	return err
}

func (p *GRPCPublisher) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
