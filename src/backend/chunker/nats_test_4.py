import asyncio
import nats
from nats.errors import TimeoutError
from nats.aio.msg import Msg
from nats.aio.client import Client as NATS
from nats.js.api import ConsumerConfig, StreamConfig, AckPolicy, DeliverPolicy, RetentionPolicy
import logging
from nats.js.errors import NotFoundError, BadRequestError

from backend.chunker.gen_types.chunking_data_pb2 import ChunkingData

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

async def main():
    nc = await nats.connect(servers=["nats://127.0.0.1:4222"])

    # Create JetStream context.
    js = nc.jetstream()

    # Create the stream configuration if it does not exist
    stream_config = StreamConfig(
        name="sample-stream",
        subjects=["foo"],
        # A work-queue retention policy satisfies a very common use case of queuing up messages that are intended to be processed once and only once.
        # https://natsbyexample.com/examples/jetstream/workqueue-stream/go
        retention=RetentionPolicy.WORK_QUEUE
    )

    try:
        await js.add_stream(stream_config)
    except BadRequestError as e:
        if e.code == 400:
            logger.info("Jetstream stream was using a different configuration. Destroying and recreating with the right configuration")
            try:
                await js.delete_stream(stream_config.name)
                await js.add_stream(stream_config)
                logger.info("Jetstream stream re-created successfully")
            except Exception as e:
                logger.exception(f"Exception while deleting and recreating Jetstream: {e}")

    # Define consumer configuration
    consumer_config = ConsumerConfig(
        durable_name="durable_chunkdata",
        ack_wait=30,  # 30 seconds
        max_deliver=3,
        ack_policy=AckPolicy.EXPLICIT,
        # DeliverPolicy.ALL is mandatory when setting  retention=RetentionPolicy.WORK_QUEUE for StreamConfig
        deliver_policy=DeliverPolicy.ALL 
    )

    # Subscribe to the subject
    try:
        await js.subscribe(subject=stream_config.subjects[0], cb=message_handler, manual_ack=True, config=consumer_config)
        logger.info("Subscribed to JetStream successfully")
    except Exception as e:
        logger.error(f"Can't subscribe to JetStream: {e}")

    # Keep the client running to listen to messages
    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        await nc.close()

async def message_handler(msg: Msg):
    try:
        logger.info("Chunking start working....")
        
        # Process the message (example code for chunking data)
        chunking_data = ChunkingData()
        chunking_data.ParseFromString(msg.data)
        # logger.info(f"URL: {chunking_data.url}")
        # logger.info(f"File Type: {chunking_data.file_type}")

        logger.info(f"Received message: {msg.data.decode()}")
        await msg.ack_sync()
        logger.info("Message acknowledged successfully")
    except Exception as e:
        logger.error(f"Chunking failed to process chunking data: {e}")
        await msg.nak()

if __name__ == '__main__':
    asyncio.run(main())
