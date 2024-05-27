import asyncio
from typing import List
from lib.milvus_db import Milvus_DB
from lib.chunked_item import ChunkedItem
from lib.spider_selenium import SeleniumSpider
from gen_types.chunking_data_pb2 import ChunkingData, FileType
from lib.chunker_base import BaseChunker
from lib.spider_bs4 import BS4Spider  # Ensure you import the BS4Spider class correctly
import logging
from datetime import datetime, timezone

class URLChunker(BaseChunker):

    async def chunk(self, data: ChunkingData):
        try:
            start_time = datetime.now(timezone.utc)
            self.logger.info(f"Starting BS4Spider URL: {data.url} at {start_time.isoformat()}")
            
            spider = BS4Spider(data.url)
            collected_data = spider.process_page(data.url)
            
            if not collected_data:
                self.logger.warning(f"BS4Spider was not able to retrieve any content for {data.url}, switching to SeleniumSpider")
                self.logger.warning("BS4Spider is disabled, shall be re-enabled and tested as it is not working 100%")
                self.logger.info(f"Starting SeleniumSpider for: {data.url}")
                spider = SeleniumSpider(data.url)
                collected_data = spider.process_page(data.url)

            if collected_data:
                milvus_db = Milvus_DB()
                # delete db from previous added chunks and vectors
                milvus_db.delete_by_document_id(document_id=data.document_id, collection_name=data.collection_name)

                # storing the new chunks in milvus
                for item in collected_data:
                    chunks = self.split_data(item.content, item.url)
                    for chunk, url in chunks:
                        milvus_db.store_chunk(content=chunk, data=data)
                        result_size_kb = len(chunk.encode('utf-8')) / 1024
                        self.logger.info(f"Chunk size for {url}: {result_size_kb:.2f} KB")
                        self.logger.info(f"{url} chunk content: {chunk}")
                        # adding some deplay not to flood milvus wit ton of requests 
                        # await asyncio.sleep(0.5)
            else:
                self.logger.warning(f"No content found for {data.url} using either BS4Spider or SeleniumSpider.")

            finish_time = datetime.now(timezone.utc)
            self.logger.info(f"Finished processing URL: {data.url} at {finish_time.isoformat()}")

            elapsed_time = finish_time - start_time
            self.logger.info(f"Elapsed time for processing {data.url}: {elapsed_time}")
            self.logger.info(f"Number of URLs analyzed: {len(collected_data)}")

            return collected_data

        except Exception as e:
            self.logger.error(f"Error: Failed to process chunking data: {e}")
            




# from typing import List
# from lib.chunked_list import ChunkedList
# from lib.spider_selenium import SeleniumSpider
# from gen_types.chunking_data_pb2 import ChunkingData, FileType
# from lib.chunker_base import BaseChunker
# from lib.spider_bs4 import BS4Spider  # Ensure you import the BS4Spider class correctly
# import logging
# from datetime import datetime, timezone

# class URLChunker(BaseChunker):

#     def chunk(self, data: ChunkingData):
#         try:
#             start_time = datetime.now(timezone.utc)
#             self.logger.info(f"Starting BS4Spider URL: {data.url} at {start_time.isoformat()}")
            
#             spider = BS4Spider(data.url)
#             # List[ChunkedList] collected_data
#             collected_data = spider.process_page(data.url)
            

#             if not collected_data:
#                 self.logger.warning(f"BS4Spider was not able to retrieve any content for {data.url}, switching to SeleniumSpider")
#                 self.logger.warning("BS4Spider is disabled, shall be re-enabled and tested as it is not working 100%")
#                 # self.logger.info(f"Starting SeleniumSpider for: {data.url}")
#                 # spider = SeleniumSpider(data.url)
#                 # collected_data = spider.process_page(data.url)

#             if collected_data:
#                 for item in collected_data:
#                     # chunk (define a method in the base class so that it can be called from here, call it split data. shall return the splitted data and the original doc or url)
#                     # store chunkes in milvus for each item coming from the process above store the content in milvusl, pass the object to MilvusDB 
#                     # the method shall check if the document has been already saved and if so upodate or delete and store again
#                     result_size_kb = len(item.content.encode('utf-8')) / 1024
#                     self.logger.info(f"Result size for {item.url}: {result_size_kb:.2f} KB")
#                     self.logger.info(f"{item.url} content: {item.content}")
#             else:
#                 self.logger.warning(f"No content found for {data.url} using either BS4Spider or SeleniumSpider.")

#             finish_time = datetime.now(timezone.utc)
#             self.logger.info(f"Finished processing URL: {data.url} at {finish_time.isoformat()}")

#             elapsed_time = finish_time - start_time
#             self.logger.info(f"Elapsed time for processing {data.url}: {elapsed_time}")
#             self.logger.info(f"Number of URLs analyzed: {len(collected_data)}")

#             return collected_data

#         except Exception as e:
#             self.logger.error(f"Error: Failed to process chunking data: {e}")
#             return []







# # from lib.spider_selenium import SeleniumSpider
# # from gen_types.chunking_data_pb2 import ChunkingData, FileType
# # from lib.chunker_base import BaseChunker
# # from lib.spider_bs4 import BS4Spider  # Ensure you import the BS4Spider class correctly
# # import logging
# # from datetime import datetime, timezone

# # class URLChunker(BaseChunker):

# #     def chunk(self, data: ChunkingData):
# #         try:
# #             self.logger.info(f"Strating BS4Spider URL: {data.url} at {datetime.now(timezone.utc).isoformat()}")
# #             spider = BS4Spider(data.url)
# #             spider.process_page(data.url)
# #             collected_data = spider.get_collected_data()


# #             if collected_data:
# #                 for item in collected_data:
# #                     result_size_kb = len(item.content.encode('utf-8')) / 1024
# #                     self.logger.info(f"Result size for {item.url}: {result_size_kb:.2f} KB")
# #                     self.logger.info(f"{item.url} content {item.content} ")
# #             else:
# #                 self.logger.warning(f"BS4Spider was not able to retieve any content for {data.url} switching to SeleniumSpider")
# #                 self.logger.warning("BS4Spider is disabled, shall be re enabled and tested it is not working 100%")
# #                 self.logger.info(f"Strating  BS4Spider for: {data.url}")
# #                 # spider = SeleniumSpider(data.url)
# #                 # spider.process_page(data.url)
# #                 # collected_data = spider.get_collected_data()

# #             self.logger.info(f"Finished BS4Spider URL: {data.url} at {datetime.now(timezone.utc).isoformat()}")
# #             return collected_data
# #         except Exception as e:
# #             self.logger.error(f"error Failed to process chunking data: {e}")
# #             return []