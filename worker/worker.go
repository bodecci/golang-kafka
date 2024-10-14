package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	topic := "comments"
	// Connect to Kafka consumer
	consumer, err := connectConsumer([]string{"localhost:29092"}, topic)
	if err != nil {
		panic(err) // Panic is used here for fatal errors
	}
	defer consumer.Close() // Ensure consumer connection is closed on exit

	fmt.Println("Consumer started")
	sigchan := make(chan os.Signal, 1)
	// Set up signal handling to gracefully shutdown the consumer
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	msgCount := 0

	// Channel to signal when processing is complete
	doneCh := make(chan struct{})

	// Goroutine for consuming messages
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println("Error received from consumer:", err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received message Count: %d: | Topic (%s) | Message (%s)\n", msgCount, string(msg.Topic), string(msg.Value))
			case <-sigchan:
				fmt.Println("Interrupt detected")
				doneCh <- struct{}{} // Signal shutdown
			}
		}
	}()

	<-doneCh // Wait for shutdown signal
	fmt.Println("Processed", msgCount, "messages")
}

// connectConsumer creates and returns a Kafka consumer connected to specified brokers and subscribed to a topic
func connectConsumer(brokersUrl []string, topic string) (sarama.PartitionConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create new consumer at the cluster level
	client, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// Consume from the specified topic and partition
	consumer, err := client.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		client.Close() // Ensure to close client if partition consumer fails to start
		return nil, fmt.Errorf("failed to start consumer for partition: %w", err)
	}

	return consumer, nil
}
