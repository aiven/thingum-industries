package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/aiven/temp/avro"
	"github.com/joho/godotenv"
	"github.com/riferrei/srclient"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	godotenv.Load()

	// set up a consumer and topics
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        os.Getenv("BROKER_URI"),
		"security.protocol":        "SSL",
		"ssl.ca.location":          "../ca.pem",
		"ssl.certificate.location": "../service.cert",
		"ssl.key.location":         "../service.key",
		"group.id":                 "CG2",
		"auto.offset.reset":        "earliest",
	})

	if err != nil {
		panic(err)
	}
	c.SubscribeTopics([]string{"factory4"}, nil)

	// wrangle the creds for the schema registry
	registry_uri := os.Getenv("REGISTRY_URI")
	registry, _ := url.Parse(registry_uri)
	reg_user := registry.User.Username()
	reg_pass, _ := registry.User.Password()
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(registry.Scheme + "://" + registry.Host)
	schemaRegistryClient.SetCredentials(reg_user, reg_pass)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			schemaID := binary.BigEndian.Uint32(msg.Value[1:5])
			schema, err := schemaRegistryClient.GetSchema(int(schemaID))
			if err != nil {
				panic(fmt.Sprintf("Error getting the schema with id '%d' %s", schemaID, err))
			}
			native, _, _ := schema.Codec().NativeFromBinary(msg.Value[5:])
			value, _ := schema.Codec().TextualFromNative(nil, native)

			// unmarshal
			var sensor avro.MachineSensor
			json.Unmarshal([]byte(value), &sensor)
			fmt.Printf("%#v\n", sensor)
		} else {
			fmt.Printf("Error consuming the message: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}
