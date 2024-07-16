# RPiDataMicroservices

Installation and Setup Instructions:

1. Install PostgreSQL and create database:

```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
sudo -u postgres psql
CREATE DATABASE sensordb;
CREATE USER username WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE sensordb TO username;
```

2. Install RabbitMQ if not already installed:

```bash
sudo apt-get install rabbitmq-server
sudo systemctl enable rabbitmq-server
sudo systemctl start rabbitmq-server
```

3. Install Go dependencies:

```bash
go mod tidy
```

4. Generate protobuf code:

```bash
make protoc
```

5. Build all services:

```bash
make all
```

6. Run the services (in separate terminals):

```bash
# Terminal 1 - Run sensor service
make run-service

# Terminal 2 - Run database service
make run-db

# Terminal 3 - Run sensor data generator
make run-go
```
