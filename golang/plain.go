package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aiven/temp/avro"
	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

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

	topic := "factory5"

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 5; i++ {
		mySensorReading := avro.MachineSensor{
			Machine: "MagicMaker4000",
			Sensor:  "oven_temp",
			Value:   (100*rand.Float32() + 150),
			Units:   "C",
		}
		reading, _ := json.Marshal(mySensorReading)

		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic: &topic, Partition: kafka.PartitionAny},
			Value: reading}, nil)
	}

	p.Flush(15 * 1000)

	fmt.Println("OK")
}
