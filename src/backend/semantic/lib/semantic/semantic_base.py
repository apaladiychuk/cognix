import os
import logging
import time
from langchain_text_splitters import RecursiveCharacterTextSplitter
from lib.gen_types.semantic_data_pb2 import SemanticData
from typing import List, Tuple
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

chunk_size = int(os.getenv('CHUNK_SIZE', 500))
chunk_overlap = int(os.getenv('CHUNK_OVERLAP', 3))
temp_path = os.getenv('LOCAL_TEMP_PATH', "../temp")


class BaseSemantic:
    def __init__(self):
        self.logger = logging.getLogger(self.__class__.__name__)

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
