package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/aiven/temp/avro"
	"github.com/joho/godotenv"
	"github.com/riferrei/srclient"
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

	topic := "factory4"

	// wrangle the creds for the schema registry
	registry_uri := os.Getenv("REGISTRY_URI")
	registry, _ := url.Parse(registry_uri)
	reg_user := registry.User.Username()
	reg_pass, _ := registry.User.Password()
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(registry.Scheme + "://" + registry.Host)
	schemaRegistryClient.SetCredentials(reg_user, reg_pass)

	// get the schema, or create it if needed
	schema, err := schemaRegistryClient.GetLatestSchema(topic, false)

	if schema == nil {
		fmt.Println("No schema, make one")
		schemaBytes, readErr := ioutil.ReadFile("../machine_sensor.avsc")
		fmt.Printf("%#v\n", readErr)
		fmt.Println(string(schemaBytes))
		schema, err = schemaRegistryClient.CreateSchema(topic, string(schemaBytes), srclient.Avro, false)
		if err != nil {
			panic(fmt.Sprintf("Error creating the schema %s", err))
		}
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 5; i++ {
		// prepare the schema ID part of the message
		schemaIDBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

		mySensorReading := avro.MachineSensor{
			Machine: "MagicMaker4000",
			Sensor:  "oven_temp",
			Value:   (100*rand.Float32() + 150),
			Units:   "C",
		}

		var valueBytesBuffer bytes.Buffer
		mySensorReading.Serialize(&valueBytesBuffer)
		valueBytes := valueBytesBuffer.Bytes()

		var recordValue []byte
		recordValue = append(recordValue, byte(0))
		recordValue = append(recordValue, schemaIDBytes...)
		recordValue = append(recordValue, valueBytes...)

		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic: &topic, Partition: kafka.PartitionAny},
			Value: recordValue}, nil)
	}

	p.Flush(15 * 1000)

	fmt.Println("OK")
}
