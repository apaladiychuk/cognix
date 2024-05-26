import logging
from gen_types.chunking_data_pb2 import ChunkingData, FileType

class BaseChunker:
    def __init__(self):
        self.logger = logging.getLogger(self.__class__.__name__)

    def chunk(self, data: ChunkingData):
        raise NotImplementedError("Chunk method needs to be implemented by subclasses")
