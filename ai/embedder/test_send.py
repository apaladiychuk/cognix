import pulsar

def main():
    try:
        # Pulsar service URL and topic to publish to
        service_url = 'pulsar://localhost:6650'
        topic = 'persistent://public/default/my-topic'

        # Create a Pulsar client
        client = pulsar.Client(service_url)

        # Create a producer on the specified topic
        producer = client.create_producer(topic)

        # Message to be sent
        message_content = 'Hello, Pulsar!'

        # Send a message
        producer.send((message_content).encode('utf-8'))
        print(f'Sent message: "{message_content}"')

        # Clean up: close the producer and client
        producer.close()
        client.close()
    except Exception as e:
        print(f"Failed to process message: {e}")
        # Message failed to process
        #consumer.negative_acknowledge(msg)

if __name__ == '__main__':
    main()
