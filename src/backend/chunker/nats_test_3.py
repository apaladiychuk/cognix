import asyncio
import nats
from nats.errors import TimeoutError
from nats.aio.msg import Msg
import logging
from nats.js.api import ConsumerConfig, StreamConfig, AckPolicy, DeliverPolicy, RetentionPolicy
from nats.js.errors import NotFoundError, BadRequestError

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

async def main():
    nc = await nats.connect(servers=["nats://127.0.0.1:4222"])
    
    # Create JetStream context.
    js = nc.jetstream()

    # Persist messages on 'foo's subject.
    # Create the stream and consumer configuration if they do not exist
    stream_config = StreamConfig(name="sample-stream", subjects=["foo"], retention=RetentionPolicy.WORK_QUEUE)
    try:
        await js.add_stream(stream_config)
    except BadRequestError as e:
        if e.err_code == 400:            
            try:
                await js.delete_stream(stream_config.name)
                logger.info("Jetstream was using a differente configuration. Destroying and recreating with the right configuration")
                await js.add_stream(stream_config)
                logger.info("Jetstream re-created succesfully") 
            except Exception as e:
                logger.exception(f"Exception while deleting and recreating Jetstream {e}")
    # await js.add_stream(name="sample-stream", subjects=["foo"])

    ########## publisher
    for i in range(0, 3):
        ack = await js.publish("foo", f"ciao mondocrudele :) :) : {i}".encode())
        print(ack)

    ######### subscriber


    # # Create pull based consumer on 'foo'.
    # psub = await js.pull_subscribe("foo", "psub")
    # msgs = await psub.fetch(1, timeout=None)
    # for msg in msgs:
    #     print(msg)
    #     await msg.ack() # <-- looks like this has no effect messages are still there and got worked again after ack
    #     #nc.publish("foo", payload=msg)
    #     await nc.flush(0.500)
    #     # si = await js.stream_info()
    #     #assertEquals(si..state.messages, 0);
    #     print(msg)
    #     print("\n\n")

    # # # Fetch and ack messagess from consumer.
    # for i in range(0, 10):
    #     msgs = await psub.fetch(1)
    #     for msg in msgs:
    #         print(msg)
    #         await msg.ack() # <-- looks like this has no effect messages are still there and got worked again after ack
            
    #         #nc.publish("foo", payload=msg)
    #         await nc.flush(0.500)
    #         # si = await js.stream_info()
    #         #assertEquals(si..state.messages, 0);
            
            
    #         print(msg)
    #         print("\n\n")

            

    # Create single ephemeral push based subscriber.
    # sub = await js.subscribe("foo")
    # msg = await sub.next_msg()
    # await msg.ack()

    # Create single push based subscriber that is durable across restarts.
    # sub = await js.subscribe("foo", durable="myapp")
    # msg = await sub.next_msg()
    # print(f"what this is doing?{msg}")
    # await msg.ack()

    # # Create deliver group that will be have load balanced messages.
    # async def qsub_a(msg):
    #     print("QSUB A:", msg)
    #     await msg.ack()

    # async def qsub_b(msg):
    #     print("QSUB B:", msg)
    #     await msg.ack()
    # await js.subscribe("foo", "workers", cb=qsub_a)
    # await js.subscribe("foo", "workers", cb=qsub_b)

    # for i in range(0, 10):
    #     ack = await js.publish("foo", f"hello world: {i}".encode())
    #     print("\t", ack)

    # Create ordered consumer with flow control and heartbeats
    # that auto resumes on failures.
    # osub = await js.subscribe("foo", ordered_consumer=True, cb=message_handler)
    
    
    
    
    # data = bytearray()

    # while True:
    #     try:
    #         msg = await osub.next_msg()
    #         data.extend(msg.data)
            
    #         msg.ack_sync()
    #         print(f"received data {msg}")
    #     except TimeoutError:
    #         break
    # print("All data in stream:", len(data))
    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        await nc.close()

async def message_handler(msg: Msg):
    try:
        logger.info("Chunking start working....")
        
        # chunking_data = ChunkingData()
        # chunking_data.ParseFromString(msg.data)
        
        # logger.info(f"URL: {chunking_data.url}")
        # logger.info(f"File Type: {chunking_data.file_type}")
        msg.ack()
        logger.info(f"message: {msg}")

        # do some work with chunking_data..

        msg.ack_sync()
        logger.info("Chunking finished working message shall be acked....")
    except Exception as e:
        logger.error(f"Chunking failed to process chunking data: {e}")
        await msg.nak()

if __name__ == '__main__':
    asyncio.run(main())