package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type doorSense struct {
	Location string `json:"location"`
	State    string `json:state"`
}

func main() {

	godotenv.Load()

	// create a producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":        os.Getenv("BROKER_URI"),
		"security.protocol":        "SSL",
		"ssl.ca.location":          "../ca.pem",
		"ssl.certificate.location": "../service.cert",
		"ssl.key.location":         "../service.key",
	})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	// event handler for messages
	go func() {
		for event := range p.Events() {
			switch ev := event.(type) {
			case *kafka.Message:
				message := ev
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Error delivering the message '%s'\n", message.Key)
				} else {
					fmt.Printf("Message delivered successfully!\n")
				}
			}
		}
	}()

	topic := "door-sensor"

	for i := 0; i < 5; i++ {
		sensor := doorSense{
			Location: "Front door",
			State:    "open",
		}
		reading, _ := json.Marshal(sensor)

		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic: &topic, Partition: kafka.PartitionAny},
			Value: reading}, nil)
	}

	p.Flush(15 * 1000)

	fmt.Println("OK")
}
