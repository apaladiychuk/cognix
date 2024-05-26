import os
import asyncio
import nats
from pyclbr import Class
from nats.aio.client import Client as NATS
from nats.aio.msg import Msg
from nats.js import JetStreamContext
from google.protobuf.json_format import Parse, MessageToJson
from google.protobuf import message as _message
from nats.js.api import ConsumerConfig, StreamConfig, AckPolicy, DeliverPolicy, RetentionPolicy
from nats.js.errors import NotFoundError, BadRequestError
from nats.js.client import JetStreamContext
from datetime import datetime
from nats.js.errors import NotFoundError
import logging
import uuid  
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Get nats url from env 
nats_url = os.getenv('NATS_URL', 'nats://127.0.0.1:4222').upper()

class JetStreamEventSubscriber:     
    def __init__(self, stream_name: str, subject: str, proto_message_type: _message.Message):
        self.stream_name = stream_name
        self.subject = subject
        self.proto_message_type = proto_message_type
        self.event_handler = None
        self.nc = NATS()
        self.js = None
        self.logger = logging.getLogger(self.__class__.__name__)

    async def connect_and_subscribe(self):
        # Connect to NATS
        await self.nc.connect(servers=[nats_url])
        # Create JetStream context
        self.js = self.nc.jetstream()

        # Create the stream configuration
        stream_config = StreamConfig(
            name=self.stream_name,
            subjects=[self.subject],
            # A work-queue retention policy satisfies a very common use case of queuing up messages that are intended to be processed once and only once.
            # https://natsbyexample.com/examples/jetstream/workqueue-stream/go
            retention=RetentionPolicy.WORK_QUEUE
            #retention=RetentionPolicy.LIMITS        
        )
        
        try:
            await self.js.add_stream(stream_config)
        except BadRequestError as e:
            if e.code == 400:
                self.logger.info("Jetstream stream was using a different configuration. Destroying and recreating with the right configuration")
                try:
                    await self.js.delete_stream(stream_config.name)
                    await self.js.add_stream(stream_config)
                    self.logger.info("Jetstream stream re-created successfully")
                except Exception as e:
                    self.logger.exception(f"Exception while deleting and recreating Jetstream: {e}")

        # # Create single ephemeral push based subscriber.
        # sub = await self.js.subscribe("foo")
        # msg = await sub.next_msg()
        # #await msg.ack()
        # self.message_handler(msg=msg)

        # Define consumer configuration
        consumer_config = ConsumerConfig(
            #name=f"consumer_name_{uuid.uuid4()}",
            #name=self.stream_name,
            #durable_name="durable_chunkdata",
            # Generate a unique durable name
            #durable_name=f"durable_{uuid.uuid4()}",  
            ack_wait=30,  # 30 seconds
            max_deliver=3,
            ack_policy=AckPolicy.EXPLICIT,
            # DeliverPolicy.ALL is mandatory when setting  retention=RetentionPolicy.WORK_QUEUE for StreamConfig
            deliver_policy=DeliverPolicy.ALL,
            #filter_subject="chunking.event.>"
        )

        # consumer_config = ConsumerConfig(
        #     ack_wait=900,
        #     max_deliver=3, 
        #     #max_ack_pending=1, 
        #     ack_policy=AckPolicy.EXPLICIT,
        #     deliver_policy=DeliverPolicy.ALL,
        # ) 


        
        # Subscribe to the subject
        try:
            
            self.js.add_consumer


            # psub = await self.js.pull_subscribe(stream=stream_config.name, subject="durable_chunkdata")
            psub = await self.js.pull_subscribe(
                subject=self.subject,
                stream=stream_config.name,
                durable="worker",
                config=consumer_config,
            )
            # psub.fetch()
            while True:
                try:
                    await asyncio.sleep(2)
                    msgs = await psub.fetch(1, timeout=5)
                    for msg in msgs:
                        # This is a life saver. Idk what it does. 
                        await msg.ack_sync()
                        await self.message_handler(msg)
                except TimeoutError:
                    print("fetch timed out . Retrying")
                    pass
            # await self.js.subscribe(subject=stream_config.subjects[0], cb=self.message_handler)#, config=consumer_config)
            self.logger.info("Subscribed to JetStream successfully")
        except Exception as e:
            self.logger.error(f"Can't subscribe to JetStream: {e}")

    async def message_handler(self, msg: Msg):
        try:        
            if self.event_handler:
                await self.event_handler(msg)
        except Exception as e:
            self.logger.exception(f"Failed to process message: {e}")

    def set_event_handler(self, event_handler):
        self.event_handler = event_handler

    async def close(self):
        await self.nc.close()

    async def flush(self):
        await self.nc.flush(0.500)