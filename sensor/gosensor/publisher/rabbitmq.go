package publisher

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
	var conn *amqp.Connection
	var err error

	// Retry connection
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"sensor_data", // exchange name
		"fanout",      // exchange type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQPublisher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *RabbitMQPublisher) PublishData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"sensor_data", // exchange
		"",            // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
}

func (p *RabbitMQPublisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
