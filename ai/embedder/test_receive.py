import pulsar

def main():
    # Pulsar service URL and topic to subscribe to
    service_url = 'pulsar://localhost:6650'
    topic = 'persistent://public/default/my-topic'

    # Create a Pulsar client
    client = pulsar.Client(service_url)

    # Create a consumer on the specified topic and subscription
    consumer = client.subscribe(topic, subscription_name='my-subscription')

    # Receive messages
    msg = consumer.receive()
    try:
        print(f"Received message: '{msg.data().decode('utf-8')}'")
        # Acknowledge successful processing of the message
        consumer.acknowledge(msg)
    except Exception as e:
        print(f"Failed to process message: {e}")
        # Message failed to process
        consumer.negative_acknowledge(msg)

    # Clean up: close the consumer and client
    consumer.close()
    client.close()

if __name__ == '__main__':
    main()
