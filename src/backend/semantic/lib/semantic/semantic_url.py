from lib.db.milvus_db import Milvus_DB
from lib.db.db_document import DocumentCRUD, Document
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic
from lib.spider.spider_bs4 import BS4Spider  # Ensure you import the BS4Spider class correctly
import time, uuid, logging, datetime

from readiness_probe import ReadinessProbe


class URLSemantic(BaseSemantic):
        def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
            try:
                self.logger.info("Analyzing URL")
                start_time = time.time()  # Record the start time
                self.logger.info(f"Starting BS4Spider URL: {data.url}")

                spider = BS4Spider(data.url)
                collected_data = spider.process_page(data.url, data.url_recursive)

                collected_items = 0
                if not collected_data:
                    self.logger.warning(f"😱 BS4Spider was not able to retrieve any content for {data.url}, switching to "
                                        f"SeleniumSpider")
                    self.logger.warning(
                        "😱 SeleniumSpider is disabled, shall be re-enabled and tested as it is not working 100%")
                    # self.logger.info(f"Starting SeleniumSpider for: {data.url}")
                    # spider = SeleniumSpider(data.url)
                    # collected_data = spider.process_page(data.url)

                chunking_session = uuid.uuid4()
                document_crud = DocumentCRUD(cockroach_url)

                if collected_data:
                    collected_items = self.store_collected_data(data=data, document_crud=document_crud,
                                                                collected_data=collected_data,
                                                                chunking_session=chunking_session,
                                                                ack_wait=ack_wait,
                                                                full_process_start_time=full_process_start_time,
                                                                split_data=True)
                else:
                    self.store_collected_data_none(data=data, document_crud=document_crud, chunking_session=chunking_session)

                self.log_end(collected_items, start_time)
                return collected_items
            except Exception as e:
                self.logger.error(f"❌ Failed to process semantic data: {e}")
