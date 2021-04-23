# Thingum Industries: Best Thingumabobs and Thingumajigs Ever!

(Please note that this repo is an imaginary example for demo purposes. You're very welcome to poke around and use anything you find here)

## Apache Kafka and the Avro Schema

You'll find an example Avro schema `machine_sensor.avsc` that is used in these examples. Avro is a binary format used in many places - this example is for using it with Kafka, both as a compression aid and to ensure that the data matches a given format.

## Python Producer for Kafka+Avro Messages

In the `python` folder you will find `sensor.py` which read the schema file, connects to Kakfa, and produces some mssages in the correct data format. To use this example yourself:

* Copy `.env/example` to `.env` and add the URLs for your Schema Registry and Apache Kafka broker. My example uses [Aiven](https://aiven.io) and therefore the schema registry is [Karapace](https://github.com/aiven/karapace).

* Download the certificates and keys to the top-level directory of the project.

* Enable auto topic creation on your Kafka instance (or create the "factory7" topic manually).

* Install the dependencies: `pip install -r requirements.txt`

* Run `python sensor.py` (requires Python3) to get some messages produced to your queue in Avro format.

* If you're using Aiven, the console has the viewer to see the produced messages, decoding the Avro. KafDrop or other clients would work well too :)

## Go Producer for Kafka+Avro Messages

In the `golang/` folder you will find the files used to talk Avro to Kafka from Go.

In order to support the correct data structure, the [gogen-avro](https://github.com/actgardner/gogen-avro
) project was used to generate structs from the `machine_sensor.avsc` schema. These are stored in the `avro/` folder and were generated with a command like:

```
gogen-avro avro machine_sensor.avsc
```

Similar to the steps for the Python setup, create a kafka instance and obtain the certificates and connection details. For a lazy quickstart, here are my scripts to do this via the `avn` commandline (run at the top level of the project):

```
avn service create --project dev-advocates -t kafka \
    -p business-4 kafka-demo \
    -c kafka.auto_create_topics_enable=true \
    -c schema_registry=true \
    -c kafka_rest=true

avn service wait kafka-demo

avn service get --json kafka-demo | jq ".service_uri"

avn service get --json kafka-demo | jq ".connection_info.schema_registry_uri"

avn service user-creds-download --username avnadmin kafka-demo
```

* Copy `.env.example` to `.env` and add the service and broker URIs to the file.

* Run with `go run main.go` to get some messages produced and sent to your Kafka topic.


## Describing Apache Kafka Payloads with AsyncAPI

The file `asyncapi.yaml` contains a description of what a consumer could expect from the produced messages. The AsyncAPI format can also read the Avro schema, so I didn't need to describe the fields twice to generate code or documentation from this file.

![screenshot of generated documentation](docs/screenshot.png)



