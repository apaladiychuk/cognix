import os
import asyncio
from nats.aio.client import Client as NATS
from google.protobuf.json_format import MessageToJson
from nats.aio.msg import Msg
from lib.chunker_helper import ChunkerHelper
from nats.js.api import ConsumerConfig, DeliverPolicy, AckPolicy
from datetime import datetime, time
from gen_types.chunking_data_pb2 import ChunkingData
from lib.jetstream_event_subscriber import JetStreamEventSubscriber
from lib.milvus_db import Milvus_DB
import logging
from dotenv import load_dotenv
import time
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

stream_name = os.getenv('CHUNKER_STREAM_NAME', 'chunker')
subject = os.getenv('CHUNKER_STREAM_SUBJECT', 'chunk_activity')


# Define the event handler function
async def chunking_event(msg: Msg):
    start_time = time.time()  # Record the start time
    try:
        logger.info("🔥 received chunking event, start working....")
        
        # Deserialize the message
        chunking_data = ChunkingData()
        chunking_data.ParseFromString(msg.data)
        logger.info(f"message: {chunking_data}")
        
        chunker_helper = ChunkerHelper()
        collecte_entities = await chunker_helper.workout_message(chunking_data)
        # if collected entities == 0 this means no data was stored in the vector db
        # we shall find a way to tell the user, most likley put the message in the dead letter

        # Acknowledge the message when done
        await msg.ack_sync()
        logger.info("👍 message acknowledged successfully")
    except Exception as e:
        logger.error(f"❌ chunking failed to process chunking data: {chunking_data} error: {e}")
        await msg.nak()
    finally:
        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")

async def main():
    logger.info("service starting")
    try:
        # subscribing to jest stream 
        subscriber = JetStreamEventSubscriber(
            stream_name=stream_name,
            subject=subject,
            proto_message_type=ChunkingData
        )

        subscriber.set_event_handler(chunking_event)
        await subscriber.connect_and_subscribe()
        logger.info("🚀 service started successfully")
    except Exception as e:
        logger.exception(f"❌ {e}")

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        await subscriber.close()

if __name__ == "__main__":
    asyncio.run(main())