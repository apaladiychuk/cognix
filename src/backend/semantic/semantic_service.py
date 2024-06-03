import asyncio
import logging
import os
import threading
import time
from dotenv import load_dotenv
from nats.aio.msg import Msg

from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_factory import SemanticFactory
from lib.db.jetstream_event_subscriber import JetStreamEventSubscriber
from lib.db.cockroach_db import DocumentCRUD
from readiness_probe import ReadinessProbe

# Load environment variables from .env file
load_dotenv()

# get log level from env
log_level_str = os.getenv('LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)
# get log format from env
log_format = os.getenv('LOG_FORMAT', '%(asctime)s - %(levelname)s - %(name)s - %(funcName)s - %(message)s')
# Configure logging
logging.basicConfig(level=log_level, format=log_format)
logger = logging.getLogger(__name__)

# loading from env env
nats_url = os.getenv('NATS_CLIENT_URL', 'nats://127.0.0.1:4222')
nats_connect_timeout = int(os.getenv('NATS_CLIENT_CONNECT_TIMEOUT', '30'))
nats_reconnect_time_wait = int(os.getenv('NATS_CLIENT_RECONNECT_TIME_WAIT', '30'))
nats_max_reconnect_attempts = int(os.getenv('NATS_CLIENT_MAX_RECONNECT_ATTEMPTS', '3'))
semantic_stream_name = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_NAME', 'semantic')
semantic_stream_subject = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_SUBJECT', 'chunk_activity')
semantic_ack_wait = int(os.getenv('NATS_CLIENT_SEMANTIC_ACK_WAIT', '3600'))  # seconds
semantic_max_deliver = int(os.getenv('NATS_CLIENT_SEMANTIC_MAX_DELIVER', '3'))

cockroach_url = os.getenv('COCKROACH_CLIENT_DATABASE_URL', 'postgres://root:123@cockroach:26257/defaultdb?sslmode=disable')

# Define the event handler function
async def chunking_event(msg: Msg):
    start_time = time.time()  # Record the start time
    try:
        logger.info("🔥 received chunking event, start working....")
        # Deserialize the message
        semantic_data = SemanticData()
        semantic_data.ParseFromString(msg.data)
        logger.info(f"message: {semantic_data}")

        chunker = SemanticFactory.create_chunker(semantic_data.file_type)

        eintites_analyzed = chunker.chunk(data= semantic_data, full_process_start_time=start_time, ack_wait=semantic_ack_wait)
        # collected_entities = await chunker.chunk( .workout_message(chunking_data=chunking_data,
        #                                                           start_time=start_time, ack_wait=semantic_ack_wait)
        # if collected entities == 0 this means no data was stored in the vector db
        # we shall find a way to tell the user, most likely put the message in the dead letter

        # Acknowledge the message when done
        await msg.ack_sync()
        logger.info("👍 message acknowledged successfully")
    except Exception as e:
        error_message = str(e) if e else "Unknown error occurred"
        logger.error(f"❌ failed to process chunking data error: {error_message}")
        if msg:  # Ensure msg is not None before awaiting
            await msg.nak()
    finally:
        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        logger.info(f"⏰⏰ total elapsed time: {elapsed_time:.2f} seconds")


async def main():
    # crud = DocumentCRUD(cockroach_url)
    #
    # # Insert a new document
    # new_doc_id = crud.insert_document(
    #     parent_id=None,
    #     connector_id=1,
    #     source_id='unique_source_id',
    #     url='http://example.com',
    #     signature='signature_example',
    #     chunking_session=uuid.uuid4(),
    #     analyzed=False,
    #     creation_date=func.now(),
    #     last_update=None
    # )
    # print(f"Inserted document ID: {new_doc_id}")
    #
    # # Select the document
    # document = crud.select_document(new_doc_id)
    # print(f"Selected document: {document}")
    #
    # # Update the document
    # crud.update_document(new_doc_id, url='http://newexample.com')
    #
    # # Delete the document
    # crud.delete_document(new_doc_id)
    # print(f"Deleted document ID: {new_doc_id}")
    #
    # # Deserialize the message
    # chunking_data = SemanticData()
    # chunking_data.ParseFromString(msg.data)
    # logger.info(f"message: {chunking_data}")



    # Start the readiness probe server in a separate thread
    readiness_probe = ReadinessProbe()
    readiness_probe_thread = threading.Thread(target=readiness_probe.start_server, daemon=True)
    readiness_probe_thread.start()

    # circuit breaker for chunking
    # if for reason nats won't be available
    # semantic will wait till nats will be up again
    while True:
        logger.info("🛠️ service starting..")
        try:
            # subscribing to jet stream
            subscriber = JetStreamEventSubscriber(
                nats_url=nats_url,
                stream_name=semantic_stream_name,
                subject=semantic_stream_subject,
                connect_timeout=nats_connect_timeout,
                reconnect_time_wait=nats_reconnect_time_wait,
                max_reconnect_attempts=nats_max_reconnect_attempts,
                ack_wait=semantic_ack_wait,
                max_deliver=semantic_max_deliver,
                proto_message_type=SemanticData
            )

            subscriber.set_event_handler(chunking_event)
            await subscriber.connect_and_subscribe()

            # todo add an event to JetStreamEventSubscriber to signal that connection has been established
            logger.info("🚀 service started successfully")

            while True:
                await asyncio.sleep(1)

        except KeyboardInterrupt:
            logger.info("🛑 Service is stopping due to keyboard interrupt")
        except Exception as e:
            logger.exception(f"💀 recovering from a fatal error: {e}. The process will restart in 5 seconds..")
            await asyncio.sleep(5)


if __name__ == "__main__":
    asyncio.run(main())
