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

# get log level from env
log_level_str = os.getenv('LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)
# get log format from env
log_format = os.getenv('LOG_FORMAT', '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
# Configure logging
logging.basicConfig(level=log_level, format=log_format)
logger = logging.getLogger(__name__)

# loading from env env
nats_url = os.getenv('NATS_CLIENT_URL', 'nats://127.0.0.1:4222')
nats_connect_timeout = int(os.getenv('NATS_CLIENT_CONNECT_TIMEOUT', '30'))
nats_reconnect_time_wait = int(os.getenv('NATS_CLIENT_RECONNECT_TIME_WAIT', '30'))
nats_max_reconnect_attempts = int(os.getenv('NATS_CLIENT_MAX_RECONNECT_ATTEMPTS', '3'))
chunker_stream_name = os.getenv('CHUNKER_STREAM_NAME','chunker')
chunker_stream_subject = os.getenv('CHUNKER_STREAM_SUBJECT','chunk_activity')
chunker_ack_wait = int(os.getenv('NATS_CLIENT_CHUNKER_ACK_WAIT', '3600')) # seconds
chunker_max_deliver = int(os.getenv('NATS_CLIENT_CHUNKER_MAX_DELIVER', '3'))


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
            nats_url = nats_url,
            stream_name=chunker_stream_name,
            subject=chunker_stream_subject,
            connect_timeout=nats_connect_timeout,
            reconnect_time_wait=nats_reconnect_time_wait,
            max_reconnect_attempts=nats_max_reconnect_attempts,
            ack_wait=chunker_ack_wait,
            max_deliver=chunker_max_deliver,
            proto_message_type=ChunkingData
        )

        subscriber.set_event_handler(chunking_event)
        await subscriber.connect_and_subscribe()

        # todo add an event to JetStreamEventSubscriber to signal that sonncetion has been established
        logger.info("🚀 service started successfully")

        try:
            while True:
                await asyncio.sleep(1)
        except KeyboardInterrupt:
            await subscriber.close()
    except Exception as e:
        logger.exception(f"❌ fatal: {e}")

if __name__ == "__main__":
    asyncio.run(main())