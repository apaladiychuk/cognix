import asyncio
import logging
from chunker.gen_types.chunking_data_pb2 import ChunkingData, FileType 
from google.protobuf.json_format import Parse, MessageToJson
from nats.errors import TimeoutError, NoRespondersError
from nats.aio.msg import Msg
from nats.aio.client import Client as NATS
from nats.js import JetStreamContext
from nats.js.api import ConsumerConfig, StreamConfig, AckPolicy, DeliverPolicy, RetentionPolicy
from nats.js.errors import NotFoundError


# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class JetStreamPublisher:
    def __init__(self, subject, stream_name):
        self.subject = subject
        self.stream_name = stream_name
        self.nc = NATS()
        self.js = None

    async def connect(self):
        await self.nc.connect(servers=["nats://127.0.0.1:4222"])
        self.js = self.nc.jetstream()

        # Ensure the stream is created
        try:
            # Try to add the stream, ignore if already exists
            await self.js.add_stream(name=self.stream_name, subjects=[self.subject])
        except Exception as e:
            logger.info(f"Stream creation error or already exists: {e}")

    async def publish(self, message):
        try:
            await self.js.publish(self.subject, message.SerializeToString())
            logger.info("Message published successfully!")
        except NoRespondersError:
            logger.error("No responders available for request")
        except TimeoutError:
            logger.error("Request to JetStream timed out")
        except Exception as e:
            logger.error(f"Failed to publish message: {e}")


        
    # Persist messages on 'foo's subject.
    # await js.add_stream(name="sample-stream", subjects=["foo"])

    ########### publisher
    # for i in range(0, 10):
    #     ack = await js.publish("foo", f"hello world: {i}".encode())
    #     print(ack)

    ########### subscriber

    async def close(self):
        await self.nc.close()

class JetStreamSubscriber:
    def __init__(self, stream_name, subject):
        self.stream_name = stream_name
        self.subject = subject
        self.event_handler = None
        self.nc = NATS()
        self.js = None

    async def connect(self):
        # Connect to NATS
        await self.nc.connect(servers=["nats://127.0.0.1:4222"])
        # Create JetStream context
        self.js = self.nc.jetstream()

        # Create the stream and consumer configuration if they do not exist
        stream_config = StreamConfig(name=self.stream_name, subjects=[self.subject], retention=RetentionPolicy.WORK_QUEUE)
        
        try:
            await self.js.add_stream(stream_config)
        except Exception as e:
            print(f"Stream already exists: {e}")
        
        consumer_config = ConsumerConfig(
            durable_name="durable_chunkdata",
            ack_wait= 30, #4 * 60 * 60,  # 4 hours in seconds
            max_deliver=3,
            ack_policy=AckPolicy.EXPLICIT,
            deliver_policy= DeliverPolicy.NEW, 
        )

        # Check if the consumer exists
        try:
            await self.js.consumer_info(stream=self.stream_name, consumer=consumer_config.name)
            print("Consumer already exists")
        except NotFoundError:
            # Create the consumer if it does not exist
            await self.js.add_consumer(stream=self.stream_name, config=consumer_config)

        await self.js.subscribe(self.subject, cb=self.message_handler)

    async def message_handler(self, msg: Msg):
        try:
            logger.info("Chunking start working....")
            
            chunking_data = ChunkingData()
            chunking_data.ParseFromString(msg.data)
            
            logger.info(f"URL: {chunking_data.url}")
            logger.info(f"File Type: {chunking_data.file_type}")

            # do some work with chunking_data..
            print(msg)
            await msg.ack()
            logger.info("Chunking finished working message shall be acked....")
        except Exception as e:
            logger.error(f"Chunking failed to process chunking data: {e}")
            await msg.nak()
        finally:
            # await self.nc.flush(0.500)
            print(msg)

    async def close(self):
        await self.nc.close()

    async def flush(self):
        await self.nc.flush(0.500)

async def main():

    # ############# publisher
    # # Instantiate the publisher
    # publisher = JetStreamPublisher(subject="chunking", stream_name="connector")

    # # Connect to NATS
    # await publisher.connect()

    # # Create a fake ChunkingData message
    # chunking_data = ChunkingData(
    #     url="https://help.collaboard.app/working-with-collaboard",
    #     site_map="",
    #     search_for_sitemap=True,
    #     document_id=993456789,
    #     file_type=FileType.URL
    # )

    # # Publish the message
    # await publisher.publish(chunking_data)

    # # Close the connection
    # await publisher.close()

    ############ subscriber
    subscriber = JetStreamSubscriber(
        stream_name="connector",
        subject="chunking"
    )

    await subscriber.connect()

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        await subscriber.close()

if __name__ == "__main__":
    asyncio.run(main())
