import asyncio
from nats.aio.client import Client as NATS
from chunker.gen_types.chunking_data_pb2 import ChunkingData, FileType 
from nats.errors import TimeoutError, NoRespondersError
from nats.js.api import ConsumerConfig, StreamConfig, AckPolicy, DeliverPolicy, RetentionPolicy
from nats.js.errors import NotFoundError, BadRequestError
import logging

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
        # Connect to NATS
        await self.nc.connect(servers=["nats://127.0.0.1:4222"])
        # Create JetStream context
        self.js = self.nc.jetstream()

        # Create the stream configuration
        stream_config = StreamConfig(
            name=self.stream_name,
            subjects=[self.subject],
            # A work-queue retention policy satisfies a very common use case of queuing up messages that are intended to be processed once and only once.
            # https://natsbyexample.com/examples/jetstream/workqueue-stream/go
            retention=RetentionPolicy.WORK_QUEUE
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
    async def create_stream(self):
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

    async def close(self):
        await self.nc.close()

async def main():
    # Instantiate the publisher
    publisher = JetStreamPublisher(subject="chunking", stream_name="connector")

    # Connect to NATS
    await publisher.connect()

    # Create a fake ChunkingData message
    chunking_data = ChunkingData(
        url="https://help.collaboard.app/working-with-collaboard",
        site_map="",
        search_for_sitemap=True,
        document_id=993456789,
        file_type=FileType.URL
    )

    # Publish the message
    await publisher.publish(chunking_data)

    # Close the connection
    await publisher.close()

if __name__ == "__main__":
    asyncio.run(main())

