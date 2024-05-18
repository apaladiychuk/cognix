import asyncio
from nats.aio.client import Client as NATS
import chunking_data_pb2  # Assuming the protobuf is compiled and named chunking_data_pb2.py
from nats.errors import TimeoutError, NoRespondersError

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
        await self.create_stream()

    async def create_stream(self):
        try:
            # Try to add the stream, ignore if already exists
            await self.js.add_stream(name=self.stream_name, subjects=[self.subject])
        except Exception as e:
            print(f"Stream creation error or already exists: {e}")

    async def publish(self, message):
        try:
            await self.js.publish(self.subject, message.SerializeToString())
            print("Message published successfully!")
        except NoRespondersError:
            print("No responders available for request")
        except TimeoutError:
            print("Request to JetStream timed out")
        except Exception as e:
            print(f"Failed to publish message: {e}")

    async def close(self):
        await self.nc.close()

async def main():
    # Instantiate the publisher
    publisher = JetStreamPublisher(subject="chunking", stream_name="connector")

    # Connect to NATS
    await publisher.connect()

    # Create a fake ChunkingData message
    chunking_data = chunking_data_pb2.ChunkingData(
        url="https://pippo.com/file.pdf",
        site_map="https://pippo.com/sitemap.xml",
        search_for_sitemap=True,
        document_id=123456789,
        file_type=chunking_data_pb2.FileType.PDF
    )

    # Publish the message
    await publisher.publish(chunking_data)

    # Close the connection
    await publisher.close()

if __name__ == "__main__":
    asyncio.run(main())
