import numpy as np
import pulsar
from pulsar.schema import JsonSchema, Record, Integer, String, Array, Float
from sentence_encoder import SentenceEncoder
from telemetry import OpenTelemetryManager

# Adapting your data structure to a JSON-compatible class
class DataSchema(Record):
    id = Integer()
    content = String()
    model = String()
    vector = Array(Float())

# Setup Pulsar client, producer, and consumer with JSON schema
client = pulsar.Client('pulsar://localhost:6650')
consumer = client.subscribe('embedd-request_v1', subscription_name='ai-embeddings_v1', schema=JsonSchema(DataSchema))
producer = client.create_producer('embedd-created_v1', schema=JsonSchema(DataSchema))

def process_message(msg):
    print(f"Received message: ID={msg.id}, Content={msg.content}")
    encoder = SentenceEncoder(msg.model)
    encoded_data = encoder.embed(msg.content)
    
    print("Encoded data:", encoded_data)
    # Directly assign the list of floats to the vector attribute
    # Convert NumPy array to a list before assigning it to the 'vector' field
    msg.vector = encoded_data.tolist() if isinstance(encoded_data, np.ndarray) else encoded_data
    
    return msg

def serve():
    try:
        while True:
            msg = consumer.receive()
            print("received message", msg)
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
