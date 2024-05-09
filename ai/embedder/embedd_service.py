import pulsar
from pulsar.schema import AvroSchema, Record, Integer, String, Array, Float

# Adapting your protobuf schema to Avro-compatible class
class DataSchema(Record):
    id = Integer()
    content = String()
    vector = Array(Float())

# Setup Pulsar client, producer, and consumer with schema
client = pulsar.Client('pulsar://localhost:6650')
consumer = client.subscribe('embedd-request', subscription_name='my-subscription', schema=AvroSchema(DataSchema))
producer = client.create_producer('output-topic', schema=AvroSchema(DataSchema))

def process_message(msg):
    print(f"Received message: ID={msg.id}, Content={msg.content}")
    return msg

def serve():
    try:
        while True:
            msg = consumer.receive()
            try:
                processed_data = process_message(msg.value())
                producer.send(processed_data)
                consumer.acknowledge(msg)
            except Exception as e:
                print("Failed to process message:", e)
                consumer.negative_acknowledge(msg)
    finally:
        client.close()

if __name__ == "__main__":
    serve()
