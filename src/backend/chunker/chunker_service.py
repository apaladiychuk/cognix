import asyncio
from nats.aio.client import Client as NATS
from google.protobuf.json_format import MessageToJson
from nats.aio.msg import Msg
from chunker.core.chunker_helper import ChunkerHelper
from nats.js.api import ConsumerConfig, DeliverPolicy, AckPolicy
from datetime import datetime
from chunker.gen_types.chunking_data_pb2 import ChunkingData
from chunker.core.jetstream_event_subscriber import JetStreamEventSubscriber
import logging

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Define the event handler function
async def chunking_event(chunking_data: ChunkingData, msg: Msg):
    try:
        logger.info("Chunking start working....")
        logger.info(f"URL: {chunking_data.url}")
        logger.info(f"File Type: {chunking_data.file_type}")

        # chunker_helper = ChunkerHelper()
        # result = chunker_helper.workout_message(chunking_data)
        # # logger.info(result)
        # if result:
        #     result_size_kb = len(result.encode('utf-8')) / 1024
        #     logger.info(f"Result size: {result_size_kb:.2f} KB")
        # else:
        #     result = ""
        #     logger.warning("Result is None")

        await msg.ack()
        logger.info("Chunking finished working....")
    except Exception as e:
        logger.error(f"Chunking failed to process chunking data: {e}")
        await msg.nak()


async def main():
    subscriber = JetStreamEventSubscriber(
        stream_name="connector",
        subject="chunking",
        proto_message_type=ChunkingData
    )

    subscriber.set_event_handler(chunking_event)
    await subscriber.connect()
    await subscriber.subscribe()

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        await subscriber.close()

if __name__ == "__main__":
    asyncio.run(main())


import asyncio
from nats.aio.client import Client as NATS
from google.protobuf.json_format import MessageToJson
from nats.aio.msg import Msg
from chunker.core.chunker_helper import ChunkerHelper
from nats.js.api import DeliverPolicy
from datetime import datetime
from chunker.gen_types.chunking_data_pb2 import ChunkingData
from chunker.core.jetstream_event_subscriber import JetStreamEventSubscriber
import logging

# Configure logging todo get from env
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Define the event handler function
async def chunking_event(chunking_data: ChunkingData, msg: Msg):
    try:
        logger.info("Chunking start working....")
        logger.info(f"URL: {chunking_data.url}")
        logger.info(f"File Type: {chunking_data.file_type}")

        # chunker_helper = ChunkerHelper()
        # result = chunker_helper.workout_message(chunking_data)
        
        # if result:
        #     # logger.info(result)
        #     result_size_kb = len(result.encode('utf-8')) / 1024
        #     logger.info(f"Result size: {result_size_kb:.2f} KB")

        logger.info("Chunking finished working....")
        await msg.ack()
    except Exception as e:
        logger.error(f"Chunking failed to process chunking data: {e}")
        await msg.nak()

async def main():
    subscriber = JetStreamEventSubscriber(
        stream_name="connector",
        subject="chunking",
        proto_message_type=ChunkingData
    )

    subscriber.set_event_handler(chunking_event)
    await subscriber.connect()
    await subscriber.subscribe()

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        await subscriber.close()

if __name__ == "__main__":
    asyncio.run(main())