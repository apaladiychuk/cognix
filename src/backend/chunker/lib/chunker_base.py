import logging
from gen_types.chunking_data_pb2 import ChunkingData, FileType
from typing import List, Tuple

class BaseChunker:
    def __init__(self):
        self.logger = logging.getLogger(self.__class__.__name__)

    def chunk(self, data: ChunkingData):
        raise NotImplementedError("Chunk method needs to be implemented by subclasses")
    


class BaseChunker:
    def __init__(self):
        self.logger = logging.getLogger(self.__class__.__name__)
    
    def split_data(self, content: str, url: str) -> List[Tuple[str, str]]:
        # This method should split the content into chunks and return a list of tuples (chunk, url)
        # For demonstration, let's split content by lines
        logging.warning("method not implemented, it is crucial to implement all the logic needed properly")
        chunks = content.split('\n')
        return [(chunk, url) for chunk in chunks if chunk]

