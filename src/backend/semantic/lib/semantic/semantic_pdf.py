from lib.db.milvus_db import Milvus_DB
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic
import time


class PDFSemantic(BaseSemantic):
    def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        try:
            start_time = time.time()  # Record the start time
            self.logger.info(f"Starting BS4Spider URL: {data.url}")

            # for pdfs, llamaparse far exceeds unstructured and pymudf is also better/faster from my experience

            document_content = "call the appropriate tool to open and extract"

            if document_content:
                milvus_db = Milvus_DB()

                # delete previous added chunks and vectors
                milvus_db.delete_by_document_id(document_id=data.document_id, collection_name=data.collection_name)

                chunks = self.split_data(document_content, data.url)
                for chunk, url in chunks:
                    milvus_db.store_chunk(content=chunk, data=data)
                    # await asyncio.sleep(0.5)
            else:
                self.logger.warning(f"😱 No content found for {data.url} ")

            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.info(f"Total elapsed time: {elapsed_time:.2f} seconds")
            #TODO: fix this
            return 0
        except Exception as e:
            self.logger.error(f"❌ Error: Failed to process chunking data: {e}")
