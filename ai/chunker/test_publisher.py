import asyncio
from nats.aio.client import Client as NATS
from nats.aio.jetstream import JetStreamContext
import chunkdata_pb2

async def publish():
    # Connect to NATS
    nc = NATS()
    await nc.connect()

    # Create JetStream context
    js = nc.jetstream()

    # Create the ChunkData message
    chunk = chunkdata_pb2.ChunkData(id="123", data=b"example data")

    # Serialize the message to a binary format
    data = chunk.SerializeToString()

    # Publish the message to a subject
    subject = "chunkdata"
    await js.publish(subject, data)

    print("Message published successfully")
    await nc.close()

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(publish())

