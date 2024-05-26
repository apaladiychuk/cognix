from gen_types.chunking_data_pb2 import ChunkingData, FileType
from lib.chunker_base import BaseChunker
from lib.spider_bs4 import BS4Spider  # Ensure you import the BS4Spider class correctly
import logging

class URLChunker(BaseChunker):

    def chunk(self, data: ChunkingData):
        try:
            self.logger.info(f"URLChunker started: {data.url}")
            spider = BS4Spider(data.url)
            spider.process_page(data.url)
            collected_data = spider.get_collected_data()


            if collected_data:
                for item in collected_data:
                    result_size_kb = len(item.content.encode('utf-8')) / 1024
                    self.logger.info(f"Result size for {item.url}: {result_size_kb:.2f} KB")
                    self.logger.info(f"{item.url} content {item.content} ")
            else:
                self.logger.warning(f"URLChunker result for {data.url} is None")


            self.logger.info(f"URLChunker finished: {data.url}")
            return collected_data
        except Exception as e:
            self.logger.error(f"URLChunker error Failed to process chunking data: {e}")
            return []

# from chunker.gen_types.chunking_data_pb2 import ChunkingData, FileType
# from chunker.core.chunker_base import BaseChunker
# from chunker.core.spider_bs4 import BS4Spider  # Ensure you import the BS4Spider class correctly

# class URLChunker(BaseChunker):
#     def chunk(self, data: ChunkingData):
#         try:
#             print(f"\n\nURLChunker started: {data.url}")
#             spider = BS4Spider(data.url)
#             spider.process_page(data.url)
#             collected_data = spider.get_collected_data()
#             print(collected_data)
#             print(f"URLChunker finished: \n\n{data.url}")
#             return collected_data
#         except Exception as e:
#             print(f"\n\nURLChunker error Failed to process chunking data: {e}")
#             return []

#     def run_chunk(self, data: ChunkingData):
#         return self.chunk(data)



# from chunker.gen_types.chunking_data_pb2 import ChunkingData, FileType
# from chunker.core.chunker_base import BaseChunker
# import requests
# from bs4 import BeautifulSoup

# class URLChunker(BaseChunker):
#     def chunk(self, data: ChunkingData):
#         try:
#             print(f"\n\nURLChunker started: {data.url}")
#             response = requests.get(data.url)
#             if response.status_code == 200:
#                 soup = BeautifulSoup(response.text, 'html.parser')
#                 collected_data = self.extract_data(soup)
#                 print(collected_data)
#                 print(f"URLChunker finished: \n\n{data.url}")
#                 return collected_data
#             else:
#                 print(f"Failed to retrieve URL: {data.url}, Status Code: {response.status_code}")
#                 return []
#         except Exception as e:
#             print(f"\n\nURLChunker error Failed to process chunking data: {e}")
#             return []

#     def extract_data(self, soup):
#         # Find all text elements within 'p', 'article', and 'div' tags
#         elements = soup.find_all(['p', 'article', 'div'])
#         paragraphs = []

#         for element in elements:
#             text = element.get_text(strip=True)
#             if text:
#                 paragraphs.append(text)

#         # Join paragraphs with two newlines for better readability
#         formatted_text = '\n\n'.join(paragraphs)
#         return formatted_text



#     def run_chunk(self, data: ChunkingData):
#         return self.chunk(data)