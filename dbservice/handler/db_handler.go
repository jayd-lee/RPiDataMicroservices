package handler

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/jayd-lee/RPiDataMicroservices/dbservice/models"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type DBHandler struct {
	db           *sql.DB
	rabbitmqConn *amqp.Connection
}

func NewDBHandler(dbConnStr, rabbitmqURL string) (*DBHandler, error) {
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create schema
	if _, err := db.Exec(models.SchemaSQL); err != nil {
		return nil, err
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}

	return &DBHandler{
		db:           db,
		rabbitmqConn: conn,
	}, nil
}

func (h *DBHandler) StartConsumingMessages() error {
	ch, err := h.rabbitmqConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the exchange (same as publisher)
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
		return err
	}

	// Declare a queue with a random name
	q, err := ch.QueueDeclare(
		"",    // name (empty = random unique name)
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(
		q.Name,        // queue name
		"",            // routing key
		"sensor_data", // exchange name
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var sensorData models.SensorData
			if err := json.Unmarshal(d.Body, &sensorData); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			if err := h.SaveSensorData(&sensorData); err != nil {
				log.Printf("Error saving sensor data: %v", err)
			} else {
				log.Printf("Successfully saved sensor data: %+v", sensorData)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func (h *DBHandler) SaveSensorData(data *models.SensorData) error {
	_, err := h.db.Exec(`
        INSERT INTO sensor_data (device_id, location, temperature, humidity, timestamp)
        VALUES ($1, $2, $3, $4, $5)
    `, data.DeviceID, data.Location, data.Temperature, data.Humidity, data.Timestamp)
	return err
}
