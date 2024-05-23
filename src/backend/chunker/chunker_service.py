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
async def chunking_event( msg: Msg):
    try:
        logger.info("Chunking start working....")
        # deserialize the message
        chunking_data = ChunkingData()
        chunking_data.ParseFromString(msg.data)
        
        logger.info(f"Received message: {chunking_data}")
        
        # chunker_helper = ChunkerHelper()
        # chunker_helper.workout_message(chunking_data)

        # await msg.ack_sync()
        logger.info("Message acknowledged successfully")
    except Exception as e:
        logger.error(f"Chunking failed to process chunking data: {chunking_data} error: {e}")
        await msg.nak()


async def main():
    subscriber = JetStreamEventSubscriber(
        stream_name="connector",
        subject="chunking",
        proto_message_type=ChunkingData
    )

    subscriber.set_event_handler(chunking_event)
    await subscriber.connect_and_subscribe()

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        await subscriber.close()

if __name__ == "__main__":
    asyncio.run(main())