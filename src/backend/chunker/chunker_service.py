import os
import asyncio
from nats.aio.client import Client as NATS
from google.protobuf.json_format import MessageToJson
from nats.aio.msg import Msg
from lib.chunker_helper import ChunkerHelper
from nats.js.api import ConsumerConfig, DeliverPolicy, AckPolicy
from datetime import datetime
from gen_types.chunking_data_pb2 import ChunkingData
from lib.jetstream_event_subscriber import JetStreamEventSubscriber
import logging
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Get log level from env 
log_level_str = os.getenv('LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)

# Get log format from env 
log_format = os.getenv('LOG_FORMAT', '%(asctime)s - %(name)s - %(levelname)s - %(message)s')

# Configure logging
logging.basicConfig(level=log_level, format=log_format)
logger = logging.getLogger(__name__)


# Define the event handler function
async def chunking_event( msg: Msg):
    try:
        logger.info("Chunking start working....")
        # deserialize the message
        chunking_data = ChunkingData()
        chunking_data.ParseFromString(msg.data)
        
        logger.info(f"Received message: {chunking_data}")
        
        chunker_helper = ChunkerHelper()
        chunker_helper.workout_message(chunking_data)

        # await msg.ack_sync()
        logger.info("Message acknowledged successfully")
    except Exception as e:
        logger.error(f"Chunking failed to process chunking data: {chunking_data} error: {e}")
        await msg.nak()


async def main():
    try:
        subscriber = JetStreamEventSubscriber(
            stream_name="connector",
            subject="chunking",
            proto_message_type=ChunkingData
        )

        subscriber.set_event_handler(chunking_event)
        await subscriber.connect_and_subscribe()
    except Exception as e:
        logger.exception(e)

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        await subscriber.close()

if __name__ == "__main__":
    asyncio.run(main())