package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
)

// Comment represents the structure for storing comment data
type Comment struct {
	Text string `form:"text" json:"text"`
}

func main() {
	app := fiber.New()          // Initialize a new Fiber application
	api := app.Group("/api/v1") // Create a route group for version 1 API

	// Define a POST route for comments
	api.Post("/comments", createComment)
	// Start the application on port 3000
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// ConnectProducer sets up a connection to the Kafka broker
func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	// Configure the producer to wait for full acknowledgment and retry up to 5 times
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	// Create and return a synchronous Kafka producer
	producer, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	return producer, nil
}

// PushCommentToQueue sends a JSON encoded comment to the specified Kafka topic
func PushCommentToQueue(topic string, message []byte) error {
	brokersUrl := []string{"localhost:29092"}
	// Establish a connection to the Kafka producer
	producer, err := ConnectProducer(brokersUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka producer: %v", err)
	}
	defer producer.Close() // Ensure the producer connection is closed on function exit

	// Construct a message and send it to Kafka
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

// createComment handles POST requests to create new comments
func createComment(c *fiber.Ctx) error {
	cmt := new(Comment)
	// Parse the JSON body into the Comment struct
	if err := c.BodyParser(cmt); err != nil {
		log.Printf("error parsing comment: %v", err)
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Invalid request data",
		})
	}
	// Serialize the Comment struct to JSON for Kafka
	cmtBytes, err := json.Marshal(cmt)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Error marshaling comment",
		})
	}
	// Send the serialized comment to the Kafka queue
	if err := PushCommentToQueue("comments", cmtBytes); err != nil {
		log.Printf("error pushing comment to Kafka: %v", err)
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	// Return a success response
	return c.JSON(&fiber.Map{
		"success": true,
		"message": "Comment pushed successfully",
		"comment": cmt,
	})
}
