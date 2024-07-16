.PHONY: all clean build run-c run-go protoc
CC=gcc
CFLAGS=-Wall -fPIC
GO=go

all: protoc build

protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/sensor/sensor.proto

build: sensor/csensor/libsensor.a sensor/csensor/sensor sensor/gosensor/sensor sensorservice/service

sensor/csensor/libsensor.a: sensor/csensor/sensor.c sensor/csensor/sensor.h
	$(CC) $(CFLAGS) -c sensor/csensor/sensor.c -o sensor/csensor/sensor.o
	ar rcs sensor/csensor/libsensor.a sensor/csensor/sensor.o

sensor/csensor/sensor: sensor/csensor/main.c sensor/csensor/sensor.h sensor/csensor/libsensor.a
	$(CC) $(CFLAGS) sensor/csensor/main.c -L./sensor/csensor -lsensor -o sensor/csensor/sensor

sensor/gosensor/sensor: sensor/gosensor/main.go sensor/csensor/libsensor.a
	cd sensor/gosensor && $(GO) build -o sensor

sensorservice/service: sensorservice/main.go
	cd sensorservice && $(GO) build -o service

dbservice/service: dbservice/main.go
	cd dbservice && $(GO) build -o service

analyticsservice/service: analyticsservice/main.go
	cd analyticsservice && $(GO) build -o service

run-db: dbservice/service
	./dbservice/service

run-analytics: analyticsservice/service
	./analyticsservice/service

clean:
	rm -f sensor/csensor/*.o sensor/csensor/*.a sensor/csensor/sensor
	rm -f sensor/gosensor/sensor
	rm -f sensorservice/service
	rm -f proto/*.pb.go

run-c: sensor/csensor/sensor
	./sensor/csensor/sensor

run-go: sensor/gosensor/sensor
	./sensor/gosensor/sensor

run-service: sensorservice/service
	./sensorservice/service
