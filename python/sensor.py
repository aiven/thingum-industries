import os
from dotenv import load_dotenv
from random import randrange

from confluent_kafka import SerializingProducer
from confluent_kafka.schema_registry import SchemaRegistryClient
from confluent_kafka.schema_registry.avro import AvroSerializer


# callback for producer response
def gotack(err, msg):
    if err is not None:
        print("Failed to deliver message: {}".format(err))
    else:
        print("OK! Topic {} partition [{}] @ offset {}"
              .format(msg.topic(), msg.partition(), msg.offset()))


# pick up the environment vars
load_dotenv()

# connect to the registry
registry_client = SchemaRegistryClient({
    'url': os.getenv('REGISTRY_URL')
})

# get the schema
with open(os.getenv('AVSC_FILENAME')) as f:
    schema = f.read()

# print(schema)

# set up the producer
producer = SerializingProducer({
    'bootstrap.servers': os.getenv('SERVICE_URL'),
    'security.protocol': 'SSL',
    'ssl.ca.location': '../ca.pem',
    'ssl.certificate.location': '../service.cert',
    'ssl.key.location': '../service.key',
    'value.serializer': AvroSerializer(schema, registry_client)
})

# send five records
for i in range(5):
    val = {
        "machine": "MagicMaker3000",
        "sensor": "tank_level",
        "value": randrange(750),
        "units": "ml"
    }

    producer.produce(topic="factory7", value=val, on_delivery=gotack)

producer.flush()
