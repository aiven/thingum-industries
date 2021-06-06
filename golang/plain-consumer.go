package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	godotenv.Load()

	// set up a consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        os.Getenv("BROKER_URI"),
		"security.protocol":        "SSL",
		"ssl.ca.location":          "../ca.pem",
		"ssl.certificate.location": "../service.cert",
		"ssl.key.location":         "../service.key",
		"group.id":                 "CG1",
		"auto.offset.reset":        "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"factory5"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}
