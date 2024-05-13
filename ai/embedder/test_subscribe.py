import asyncio
from nats.aio.client import Client as NATS
from nats.js.api import JetStreamContext
import embedd_data_pb2  # Import the generated protobuf module

async def handle_task(msg):
    # Deserialize the Task message
    task = embedd_data_pb2.Task()
    task.ParseFromString(msg.data)
    print(f"Received task to: {task.content}")
    # Acknowledge the message
    await msg.ack()

async def run():
    # Connect to NATS
    nc = NATS()
    await nc.connect(servers=["nats://localhost:4222"])
    js = nc.jetstream()

    # Subscribe to tasks
    await js.subscribe("tasks.subject", cb=handle_task)

    # Keep the connection alive to continue receiving tasks
    await asyncio.Event().wait()

if __name__ == '__main__':
    asyncio.run(run())
