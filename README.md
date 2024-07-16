# RPiDataMicroservices

Installation and Setup Instructions:

1. Install RabbitMQ if not already installed:

```bash
sudo apt-get install rabbitmq-server
sudo systemctl enable rabbitmq-server
sudo systemctl start rabbitmq-server
```

2. Install Go dependencies:

```bash
go mod tidy
```

3. Generate protobuf code:

```bash
make protoc
```

4. Build all services:

```bash
make all
```

5. Run the services (in separate terminals):

```bash
# Terminal 1 - Run sensor service
make run-service

# Terminal 2 - Run sensor data generator
make run-go
```
