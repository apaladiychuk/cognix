from chunker.gen_types.chunking_data_pb2 import ChunkingData, FileType
from chunker.core.chunker_base import BaseChunker
import logging

class PDFChunker(BaseChunker):
    # def __init__(self):
    #     super().__init__()

    def chunk(self, data: ChunkingData):
        try:
            self.logger.info(f"PDFChunker not implemented started: {data.url}")
            # Implement PDF chunking logic here
            self.logger.info(f"PDFChunker finished: {data.url}")
        except Exception as e:
            self.logger.error(f"PDFChunker error Failed to process chunking data: {e}")
