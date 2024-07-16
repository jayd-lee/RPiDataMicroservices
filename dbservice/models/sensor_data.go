package models

import (
	"time"
)

type SensorData struct {
	ID          int64     `db:"id"`
	DeviceID    string    `db:"device_id"`
	Location    string    `db:"location"`
	Temperature float64   `db:"temperature"`
	Humidity    float64   `db:"humidity"`
	Timestamp   time.Time `db:"timestamp"`
}

const SchemaSQL = `
CREATE TABLE IF NOT EXISTS sensor_data (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    temperature DOUBLE PRECISION NOT NULL,
    humidity DOUBLE PRECISION NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_device_timestamp ON sensor_data(device_id, timestamp);
`
