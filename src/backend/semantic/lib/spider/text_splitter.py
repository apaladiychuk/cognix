import os
from langchain_text_splitters import RecursiveCharacterTextSplitter
import logging
from typing import List, Tuple

class TextSplitter:
    def __init__(self):
        self.chunk_size = int(os.getenv('CHUNK_SIZE', 500))
        self.chunk_overlap = int(os.getenv('CHUNK_OVERLAP', 3))

    def split_data(self, content: str, url: str) -> List[Tuple[str, str]]:
        logging.warning("😱 split_data shall implement various chunk techniques and compare them")

        # Initialize the text splitter with custom parameters
        custom_text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=self.chunk_size,
            chunk_overlap=self.chunk_overlap,
            length_function=len,
            separators=['\n']
        )

        # Create the chunks
        texts = custom_text_splitter.create_documents([content])

        if texts:
            logging.info(f"created {len(texts)} chunks for {url}")
        else:
            logging.info(f"no chunk created for {url}")

        return [(chunk.page_content, url) for chunk in texts if chunk]
