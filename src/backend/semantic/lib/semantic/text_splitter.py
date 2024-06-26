import os
from langchain_text_splitters import RecursiveCharacterTextSplitter
import logging
from typing import List, Tuple

from lib.spider.chunked_item import ChunkedItem

class TextSplitter:
    chunk_size = int(os.getenv('CHUNK_SIZE', 500))
    chunk_overlap = int(os.getenv('CHUNK_OVERLAP', 3))

    @classmethod
    def create_chunked_items(cls, content: str, url: str, document_id: int, parent_id: int) -> List['ChunkedItem']:
        chunked_items = []
        custom_text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=cls.chunk_size,
            chunk_overlap=cls.chunk_overlap,
            length_function=len,
            separators=['\n']
        )

        texts = custom_text_splitter.create_documents([content])

        if texts:
            logging.info(f"created {len(texts)} chunks for {url}")
        else:
            logging.info(f"no chunk created for {url}")

        for chunk in texts:
            if chunk:
                chunked_items.append(ChunkedItem(content=chunk.page_content, url=url, document_id=document_id, parent_id=parent_id))

        return chunked_items

