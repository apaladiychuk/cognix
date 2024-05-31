from lib.milvus_db import Milvus_DB
from gen_types.chunking_data_pb2 import ChunkingData, FileType
from lib.chunker_base import BaseChunker
from pathlib import Path
import logging, time
import pymupdf

class TXTChunker(BaseChunker):
    async def chunk(self, data: ChunkingData) -> int:
        try:
            start_time = time.time()  # Record the start time
            self.logger.info(f"Starting TXTChunker for: {data.url}")
            
            doc = pymupdf.open("PyMuPDF.pdf") # open a supported document
            page = doc[0] # load the required page (0-based index)
            text = page.get_text() # extract plain text
            print(text) # process or print it:
  
            all_text = ""
            for page in doc:
                all_text += page.get_text() + chr(12)
                
                # or, with the even faster list comprehension:
                all_text = chr(12).join([page.get_text() for page in doc])


            pages = 0
            # TODO: how to read a big text file 
            # https://medium.com/@tubelwj/how-to-read-extremely-large-text-files-in-python-cddc7dbce9fc
            # atm we are doing the dumb way :) 
            document_content = Path(data.url).read_text()

            if document_content:
                pages = 1
                milvus_db = Milvus_DB()

                # delete previous added chunks and vectors
                milvus_db.delete_by_document_id(document_id=data.document_id, collection_name=data.collection_name)

                chunks = self.split_data(document_content, data.url)
                for chunk in chunks:
                    milvus_db.store_chunk(content=chunk, data=data)
                    # await asyncio.sleep(0.5)
            else:
                self.logger.warning(f"😱 No content found for {data.url} ")

            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.info(f"Total elapsed time: {elapsed_time:.2f} seconds")
            return pages
        except Exception as e:
            self.logger.error(f"❌ Error: Failed to process chunking data: {e}")