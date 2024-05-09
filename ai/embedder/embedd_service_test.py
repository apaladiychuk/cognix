import pulsar
from pulsar.schema import JsonSchema, Record, Integer, String, Array, Float

class DataSchema(Record):
    id = Integer()
    content = String()
    model = String()
    vector = Array(Float())

# Setup the Pulsar client and producer with JSON schema
client = pulsar.Client('pulsar://localhost:6650')
producer = client.create_producer('embedd-request_v1', schema=JsonSchema(DataSchema))

def send_message():
    # Prompt user for input to embed
    content_to_embed = input("Type the content you want to embed: ")
    model_to_use = "sentence-transformers/paraphrase-multilingual-mpnet-base-v2"

    # Create a new message object using the DataSchema
    message = DataSchema(id=123, content=content_to_embed, model=model_to_use, vector=[1.0, 2.0, 3.0])

    # Send the message to the topic through the producer
    producer.send(message)
    print(f"Sent message: ID={message.id}, Content={message.content}")

if __name__ == "__main__":
    try:
        send_message()
    finally:
        client.close()
