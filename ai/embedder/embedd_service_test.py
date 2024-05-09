import pulsar
from pulsar.schema import AvroSchema, Record, Integer, String, Array, Float

# Define the data schema using Pulsar's AvroSchema
class DataSchema(Record):
    id = Integer()
    content = String()
    vector = Array(Float())

# Setup the Pulsar client and producer with schema
client = pulsar.Client('pulsar://localhost:6650')
producer = client.create_producer('embedd-request', schema=AvroSchema(DataSchema))

def send_message():
    # Create a new message using the DataSchema
    message = DataSchema(id=123, content="Hello, Pulsar!", vector=[1.0, 2.0, 3.0])
    # Send the message to the topic through the producer
    producer.send(message)
    print(f"Sent message: ID={message.id}, Content={message.content}")

if __name__ == "__main__":
    send_message()
    client.close()
