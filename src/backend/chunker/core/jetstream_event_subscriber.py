import asyncio
from nats.aio.client import Client as NATS
from nats.aio.msg import Msg
from nats.js import JetStreamContext
from google.protobuf.json_format import Parse, MessageToJson
from nats.js.api import ConsumerConfig, StreamConfig, AckPolicy, DeliverPolicy
from datetime import datetime
from nats.js.errors import NotFoundError

class JetStreamEventSubscriber:
    def __init__(self, stream_name, subject, proto_message_type):
        self.stream_name = stream_name
        self.subject = subject
        self.proto_message_type = proto_message_type
        self.event_handler = None
        self.nc = NATS()
        self.js = None

    async def connect(self):
        # Connect to NATS
        await self.nc.connect(servers=["nats://127.0.0.1:4222"])
        # Create JetStream context
        self.js = self.nc.jetstream()

    async def subscribe(self):
        # Create the stream and consumer configuration if they do not exist
        stream_config = StreamConfig(name=self.stream_name, subjects=[self.subject])
        
        try:
            await self.add_stream(stream_config)
        except Exception as e:
            print(f"Stream already exists: {e}")
        
        consumer_config = ConsumerConfig(
            durable_name="durable_chunkdata",
            ack_wait= 30, #4 * 60 * 60,  # 4 hours in seconds
            max_deliver=3,
            ack_policy=AckPolicy.EXPLICIT,
            deliver_policy= DeliverPolicy.NEW,
        )

        # for dev purposes only
        # Delete the existing consumer if it exists
        # try:
        #     await self.js.delete_consumer("connector", "durable_chunkdata")
        #     print("Deleted existing consumer")
        # except NotFoundError:
        #     print("Consumer does not exist, no need to delete")

        # Check if the consumer exists
        try:
            await self.js.consumer_info(stream=self.stream_name, consumer="durable_chunkdata")
            print("Consumer already exists")
        except NotFoundError:
            # Create the consumer if it does not exist
            await self.js.add_consumer(stream=self.stream_name, config=consumer_config)

        await self.js.subscribe(self.subject, cb=self.message_handler)

    async def message_handler(self, msg):
        try:
            proto_message = self.proto_message_type()
            proto_message.ParseFromString(msg.data)
            
            if self.event_handler:
                await self.event_handler(proto_message, msg)
        except Exception as e:
            print(f"Failed to process message: {e}")

    def set_event_handler(self, event_handler):
        self.event_handler = event_handler

    async def close(self):
        await self.nc.close()

    async def flush(self):
        await self.nc.flush(0.500)