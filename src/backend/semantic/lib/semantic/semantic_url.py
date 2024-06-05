from lib.db.milvus_db import Milvus_DB
from lib.db.db_document import DocumentCRUD, Document
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic
from lib.spider.spider_bs4 import BS4Spider  # Ensure you import the BS4Spider class correctly
import time, uuid, logging, datetime

from readiness_probe import ReadinessProbe


class URLSemantic(BaseSemantic):
    def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        start_time = time.time()  # Record the start time
        self.logger.info(f"Starting BS4Spider URL: {data.url}")

        spider = BS4Spider(data.url)
        collected_data = spider.process_page(data.url, data.url_recursive)
        collected_items = 0

        if not collected_data:
            self.logger.warning(f"😱 BS4Spider was not able to retrieve any content for {data.url}, switching to "
                                f"SeleniumSpider")
            self.logger.warning("😱 SeleniumSpider is disabled, shall be re-enabled and tested as it is not working 100%")
            # self.logger.info(f"Starting SeleniumSpider for: {data.url}")
            # spider = SeleniumSpider(data.url)
            # collected_data = spider.process_page(data.url)

        chunking_session = uuid.uuid4()
        document_crud = DocumentCRUD(cockroach_url)

        if collected_data:
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
            doc.chunking_session=chunking_session
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
        else:
            # storing in the db the item setting analyzed = false cause we were not able to extract any text out of it
            # there will be no trace of it in milvus
            doc = Document(parent_id=data.document_id, connector_id=data.connector_id, source_id=data.url,
                           url=data.url, chunking_session=chunking_session, analyzed=False,
                           creation_date=datetime.datetime.utcnow(), last_update=datetime.datetime.utcnow())
            document_crud.update_document(doc)
            self.logger.warning(f"😱 no content found for {data.url} using either BS4Spider or SeleniumSpider.")

        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        self.logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")
        self.logger.info(f"📖 number of URLs analyzed: {collected_items}")
        return collected_items

            
