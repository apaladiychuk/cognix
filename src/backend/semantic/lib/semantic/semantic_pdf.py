import time
import uuid
from typing import List, Dict, Optional
from urllib.parse import urlparse, parse_qs

import pymupdf4llm

from lib.db.db_document import DocumentCRUD
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic
from lib.spider.chunked_item import ChunkedItem
from test_mistune import MarkdownSectionExtractor


class PDFSemantic(BaseSemantic):
    def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        try:
            start_time = time.time()  # Record the start time
            self.logger.info(f"Starting pdf analysis for {data.url}")
            t0 = time.perf_counter()

            markdown_content = pymupdf4llm.to_markdown(data.url)
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
