import asyncio
from nats.aio.client import Client as NATS
from google.protobuf.json_format import MessageToJson
from nats.aio.msg import Msg
from chunker.core.chunker_helper import ChunkerHelper
from nats.js.api import DeliverPolicy
from datetime import datetime
from chunker.gen_types.chunking_data_pb2 import ChunkingData
from chunker.core.jetstream_event_subscriber import JetStreamEventSubscriber


# Define the event handler function
async def chunking_event(chunking_data: ChunkingData, msg: Msg):
    try:
        print("\n\nchunking start working....")
        print(f"URL: {chunking_data.url}")
        print(f"File Type: {chunking_data.file_type}")

        chunker_helper = ChunkerHelper()
        result = chunker_helper.workout_message(chunking_data)
        print(result)

        print("chunking finished working.... \n\n")
        await msg.ack()
    except Exception as e:
        print(f"chunking failed to process chunking data: {e}")
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



# import asyncio
# from nats.aio.client import Client as NATS
# from google.protobuf.json_format import MessageToJson
# from nats.aio.msg import Msg
# from chunker.core.chunker_helper import ChunkerHelper
# from chunker.core.chunker_helper import ChunkerHelper
# #import backend.chunker.chunking_data_pb2
# from nats.js.api import DeliverPolicy
# from datetime import datetime
# from chunker.gen_types.chunking_data_pb2 import ChunkingData
# from chunker.core.jetstream_event_subscriber import JetStreamEventSubscriber


# # Define the event handler function
# async def chunking_event(chunking_data: ChunkingData, msg: Msg):
#     try:
#         # # Fake implementation: Deserialize and print the message
#         # print("Deserialized ChunkingData message:")
#         # print(datetime.now())
#         # print(f"URL: {chunking_data.url}")
#         # print(f"Site Map: {chunking_data.site_map}")
#         # print(f"Search for Sitemap: {chunking_data.search_for_sitemap}")
#         # print(f"Document ID: {chunking_data.document_id}")
#         # print(f"File Type: {chunking_data.file_type}")

#         # # Optionally, print the entire message as JSON
#         # print(f"Received data as JSON: {MessageToJson(chunking_data)}")

#         print("\n \n chunking start working....")
#         print(f"URL: {chunking_data.url}")
#         print(f"File Type: {chunking_data.file_type}")

#         # creating an instance of ChunkerHelper
#         # to start the work 
#         chunker_helper = ChunkerHelper()
#         chunker_helper.workout_message(chunking_data)
        

#         # Acknowledge the message
#         await msg.ack()
#         print("chunking finished working.... \n \n")
#     except Exception as e:
#         print(f" chunking failed to process chunking data: {e}")
#         # Optionally, do not acknowledge the message (it will be retried)
#         await msg.nak()  # Uncomment if you want to explicitly not acknowledge the message

# async def main():

#     subscriber = JetStreamEventSubscriber(
#         stream_name="connector",
#         subject="chunking",
#         proto_message_type=ChunkingData
#     )

#     # Set the event handler
#     subscriber.set_event_handler(chunking_event)

#     await subscriber.connect()
#     await subscriber.subscribe()

#     # Keep the script running to listen for messages
#     try:
#         while True:
#             await asyncio.sleep(1)
#     except KeyboardInterrupt:
#         await subscriber.close()

# if __name__ == "__main__":
#     asyncio.run(main())
