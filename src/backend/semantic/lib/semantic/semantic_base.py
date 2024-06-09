import datetime
import os
import logging
import time
import uuid

import pymupdf4llm
from langchain_text_splitters import RecursiveCharacterTextSplitter
from minio import Minio, S3Error

from lib.db.db_document import Document, DocumentCRUD
from lib.db.milvus_db import Milvus_DB
from lib.gen_types.semantic_data_pb2 import SemanticData
from typing import List, Tuple
from dotenv import load_dotenv

from lib.semantic.markdown_extractor import MarkdownSectionExtractor
from lib.spider.chunked_item import ChunkedItem
from readiness_probe import ReadinessProbe

from typing import List

# Load environment variables from .env file
load_dotenv()

chunk_size = int(os.getenv('CHUNK_SIZE', 500))
chunk_overlap = int(os.getenv('CHUNK_OVERLAP', 3))
temp_path = os.getenv('LOCAL_TEMP_PATH', "../temp")

minio_endpoint = os.getenv('MINIO_ENDPOINT', "minio:9000")
minio_access_key = os.getenv('MINIO_ACCESS_KEY', "minioadmin")
minio_secret_key = os.getenv('MINIO_SECRET_ACCESS_KEY', "minioadmin")
minio_use_ssl = os.getenv('MINIO_USE_SSL', 'false').lower() == 'true'


class BaseSemantic:
    def __init__(self):
        self.logger = logging.getLogger(self.__class__.__name__)
        self.temp_path = temp_path
        self.minio_endpoint = minio_endpoint
        self.minio_access_key = minio_access_key
        self.minio_secret_key = minio_secret_key
        self.minio_use_ssl = minio_use_ssl

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

    def convert_to_markdown(self, file_path: str) -> str:
        return pymupdf4llm.to_markdown(file_path)

    def download_from_minio(self, url: str) -> str:
        """
        Download a file from MinIO using the provided URL and save it to the specified local temporary path.

        :param url: The MinIO URL of the file to be downloaded.
        :return: The full path to the downloaded file.
        """
        # Extract bucket name and object name from the URL
        parts = url.split(':')
        bucket_name = parts[1]
        object_name = parts[-1]

        # Extract the file name from the object name
        file_name = object_name.split('-')[-1]
        # Combine the temporary path and the file name
        save_path = os.path.join(self.temp_path, file_name)

        # Initialize the MinIO client
        client = Minio(
            self.minio_endpoint,
            access_key=self.minio_access_key,
            secret_key=self.minio_secret_key,
            secure=self.minio_use_ssl  # Use SSL if minio_use_ssl is true
        )

        # Download the file from the bucket
        client.fget_object(bucket_name, object_name, save_path)
        print(f"File downloaded successfully and saved to {save_path}")
        return save_path

    def analyze_doc(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        try:
            start_time = time.time()  # Record the start time
            self.logger.info(f"Starting pdf analysis for {data.url}")
            t0 = time.perf_counter()

            #minio:tenant-c20a9f75-a363-40ea-86ef-eabcedbac7df:86eafa6e-95b3-4b29-859a-c38bd4552f26-Document.docx

            downloaded_file_path = self.download_from_minio(data.url)

            markdown_content = pymupdf4llm.to_markdown(downloaded_file_path)
            # print(markdown_content)

            extractor = MarkdownSectionExtractor()
            # sections = extractor.extract_sections(markdown_content)
            results = extractor.extract_chunks(markdown_content)
            collected_data = ChunkedItem.create_chunked_items(results, data.url)

            collected_items = 0
            if not collected_data:
                self.logger.warning(f"😱no content found in {data.url}")

            chunking_session = uuid.uuid4()
            document_crud = DocumentCRUD(cockroach_url)

            if collected_data:
                collected_items = self.store_collected_data(data=data, document_crud=document_crud,
                                                            collected_data=collected_data,
                                                            chunking_session=chunking_session,
                                                            ack_wait=ack_wait,
                                                            full_process_start_time=full_process_start_time)
            else:
                self.store_collected_data_none(data=data, document_crud=document_crud,
                                               chunking_session=chunking_session)

            self.log_end(collected_items, start_time)
            return collected_items
        except Exception as e:
            self.logger.error(f"❌ Failed to process semantic data: {e}")
