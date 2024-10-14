# Golang-Kafka

## Message Processing System

#### This application is a simple system that allows you to add a message through a REST API, then produce it in Apache Kafka and process it in a consumer (worker). The consumer performs time-consuming tasks, such as saving data to a datastore, performing aggregations, or running computations.

#### This is an Apache with Go. an exercise in using Kafka with Go, making use of the client Sarama library.

### Setup

1. Install Docker and Docker Compose

2. Run `docker-compose up -d` in the root directory. This will start Apache Kafka and Zookeeper

3. Install dependencies

4. Run `go run producer/producer.go` to start the producer which is a REST API listening on port 3000

5. Run `go run worker/worker.go` to start the consumer

### Test it out!

Send a POST request to `localhost:3000`

This will produce a message in Apache Kafka and the consumer will process it.

```bash
curl --location --request POST '0.0.0.0:3000/api/v1/comments' \
--header 'Content-Type: application/json' \
--data-raw '{ "text":"message 1" }'

curl --location --request POST '0.0.0.0:3000/api/v1/comments' \
--header 'Content-Type: application/json' \
--data-raw '{ "text":"message 2" }'
```