import os
import asyncio
from nats.aio.client import Client as NATS
from nats.aio.jetstream import JetStreamContext, Msg
import chunkdata_pb2

async def process_message(msg: Msg):
    #pulsarURL = os.environ.get['PUlSAR_URL']

    # Deserialize the message
    chunk = chunkdata_pb2.ChunkData()
    chunk.ParseFromString(msg.data)
    print(f"Received ChunkData: ID={chunk.id}, Data={chunk.data}")

    # Simulate message processing
    try:
        if chunk.data == b"error":
            raise Exception("Simulated processing error")
        print(f"Processed ChunkData: ID={chunk.id}, Data={chunk.data}")
        await msg.ack()
    except Exception as e:
        print(f"Error processing message: {e}")
        # Do not acknowledge the message to trigger a retry

async def subscribe():
    # Connect to NATS
    nc = NATS()
    await nc.connect()

    # Create JetStream context
    js = nc.jetstream()

    # Create the stream and consumer configuration if they do not exist
    await js.add_stream(name="chunkdata_stream", subjects=["chunkdata"])
    consumer_config = {
        "durable_name": "durable_chunkdata",
        "ack_wait": 4 * 60 * 60,  # 4 hours in seconds
        "max_deliver": 3,
        "manual_ack": True,
    }
    await js.add_consumer("chunkdata_stream", consumer_config)

    # Subscribe to the subject with the durable consumer
    await js.subscribe("chunkdata", "durable_chunkdata", cb=process_message)

    # Keep the subscriber running
    await asyncio.Event().wait()

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(subscribe())
