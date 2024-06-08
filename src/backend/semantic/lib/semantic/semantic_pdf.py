import time
import uuid
from typing import List, Dict, Optional
from urllib.parse import urlparse, parse_qs

from lib.db.db_document import DocumentCRUD
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic
from lib.spider.chunked_item import ChunkedItem

class PDFSemantic(BaseSemantic):
    def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        try:
            start_time = time.time()  # Record the start time
            self.logger.info(f"Starting BS4Spider URL: {data.url}")

            # for pdfs, llamaparse far exceeds unstructured and pymudf is also better/faster from my experience

            content = ""

            collected_items = 0
            chunking_session = uuid.uuid4()
            document_crud = DocumentCRUD(cockroach_url)

            if content:
                collected_data = [ChunkedItem(url=data.url, content=content)]

                collected_items = self.store_collected_data(data=data, document_crud=document_crud,
                                                            collected_data=collected_data,
                                                            chunking_session=chunking_session,
                                                            ack_wait=ack_wait,
                                                            full_process_start_time=full_process_start_time)
                self.logger.debug(f"transcript \n {content}")
            else:
                self.store_collected_data_none(data=data, document_crud=document_crud,
                                               chunking_session=chunking_session)

            self.log_end(collected_items, start_time)
            return collected_items
            # (if transcript: 1 else: 0)
        except Exception as e:
            self.logger.error(f"❌ Failed to process semantic data: {e}")
