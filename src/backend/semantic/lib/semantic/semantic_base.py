import datetime
import os
import logging
import time
import uuid

from langchain_text_splitters import RecursiveCharacterTextSplitter

from lib.db.db_document import Document, DocumentCRUD
from lib.db.milvus_db import Milvus_DB
from lib.gen_types.semantic_data_pb2 import SemanticData
from typing import List, Tuple
from dotenv import load_dotenv

from lib.spider.chunked_item import ChunkedItem
from readiness_probe import ReadinessProbe

# Load environment variables from .env file
load_dotenv()

chunk_size = int(os.getenv('CHUNK_SIZE', 500))
chunk_overlap = int(os.getenv('CHUNK_OVERLAP', 3))
temp_path = os.getenv('LOCAL_TEMP_PATH', "../temp")


class BaseSemantic:
    def __init__(self):
        self.logger = logging.getLogger(self.__class__.__name__)
        self.temp_path = temp_path

    def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        raise NotImplementedError("Chunk method needs to be implemented by subclasses")
    def keep_processing(self, full_process_start_time: float, ack_wait: int) -> bool:
        # it returns true if the difference between start_time and now is less than ack_wait
        # it returns false if the difference between start_time and now is equal or greater than ack_wait
        end_time = time.time()  # Record the end time
        elapsed_time = end_time - full_process_start_time
        return elapsed_time < ack_wait

    def split_data(self, content: str, url: str) -> List[Tuple[str, str]]:
        # This method should split the content into chunks and return a list of tuples (chunk, url)
        # For demonstration, let's split content by lines
        logging.warning("😱 split_data shall implement various chunk techniques and compare them")

        # Initialize the text splitter with custom parameters
        custom_text_splitter = RecursiveCharacterTextSplitter(
            # Set custom chunk size
            chunk_size=chunk_size,
            chunk_overlap=chunk_overlap,
            # Use length of the text as the size measure
            length_function=len,
            # Use only "\n\n" as the separator
            separators=['\n']
        )

        # Create the chunks
        texts = custom_text_splitter.create_documents([content])

        if texts:
            self.logger.info(f"created {len(texts)} chunks for {url}")
        else:
            self.logger.info(f"no chunk created for {url}")

        return [(chunk.page_content, url) for chunk in texts if chunk]

    def store_collected_data(self, data: SemanticData, document_crud: DocumentCRUD, collected_data: list[ChunkedItem],
                             chunking_session: uuid, ack_wait: int, full_process_start_time: float):
        collected_items = 0
        # verifies if the method is taking longer than ack_wait
        # if so we have to stop
        if not self.keep_processing(full_process_start_time=full_process_start_time, ack_wait=ack_wait):
            raise Exception(f"exceeded maximum processing time defined in NATS_CLIENT_SEMANTIC_ACK_WAIT of {ack_wait}")
        if self.logger.level == logging.DEBUG:
            collected_items = len(collected_data)
            self.logger.debug(f"collected {collected_items} URLs")
        milvus_db = Milvus_DB()
        # delete previous added chunks and vectors
        milvus_db.delete_by_document_id(document_id=data.document_id, collection_name=data.collection_name)
        # delete previous added documents
        document_crud.delete_by_parent_id(data.document_id)
        doc = document_crud.select_document(data.document_id)
        doc.chunking_session = chunking_session
        doc.analyzed = True
        doc.last_update = datetime.datetime.utcnow()
        document_crud.update_document_object(doc)
        # Now in this mess find the parent document!!!!
        # all children can be added randomly
        # storing the new chunks in milvus
        for item in collected_data:
            # verifies if the method is taking longer than ack_wait
            # if so we have to stop
            if not self.keep_processing(full_process_start_time=full_process_start_time, ack_wait=ack_wait):
                raise Exception(
                    f"exceeded maximum processing time defined in NATS_CLIENT_SEMANTIC_ACK_WAIT of {ack_wait}")

            # insert in milvus
            chunks = self.split_data(item.content, item.url)
            for chunk, url in chunks:
                # notifying the readiness probe that the service is alive
                ReadinessProbe().update_last_seen()

                # verifies if the method is taking longer than ack_wait
                # if so we have to stop
                if not self.keep_processing(full_process_start_time=full_process_start_time, ack_wait=ack_wait):
                    raise Exception(
                        f"exceeded maximum processing time defined in NATS_CLIENT_SEMANTIC_ACK_WAIT of {ack_wait}")

                # and finally the real job!!!
                milvus_db.store_chunk(content=chunk, data=data)

                if self.logger.level == logging.DEBUG:
                    result_size_kb = len(chunk.encode('utf-8')) / 1024
                    self.logger.debug(f"Chunk size for {url}: {result_size_kb:.2f} KB")
                    self.logger.debug(f"{url} chunk content: {chunk}")

            # let's store the doc in the relational db
            doc = Document(parent_id=data.document_id, connector_id=data.connector_id, source_id=item.url,
                           url=item.url, chunking_session=chunking_session, analyzed=True,
                           creation_date=datetime.datetime.utcnow(), last_update=datetime.datetime.utcnow())

            document_crud.insert_document_object(doc)
            collected_items += len(chunks)
        return collected_items

    def store_collected_data_none(self, data: SemanticData, document_crud: DocumentCRUD, chunking_session: uuid):
        # storing in the db the item setting analyzed = false because we were not able to extract any text out of it
        # there will be no trace of it in milvus
        doc = Document(parent_id=data.document_id, connector_id=data.connector_id, source_id=data.url,
                       url=data.url, chunking_session=chunking_session, analyzed=False,
                       creation_date=datetime.datetime.utcnow(), last_update=datetime.datetime.utcnow())
        document_crud.update_document_object(doc)
        self.logger.warning(f"😱 no content found for {data.url} using either BS4Spider or SeleniumSpider.")

    def log_end(self, collected_items, start_time):
        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        self.logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")
        self.logger.info(f"📖 number of URLs analyzed: {collected_items}")