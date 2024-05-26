from lib.spider_selenium import SeleniumSpider
from gen_types.chunking_data_pb2 import ChunkingData, FileType
from lib.chunker_base import BaseChunker
from lib.spider_bs4 import BS4Spider  # Ensure you import the BS4Spider class correctly
import logging
from datetime import datetime, timezone

class URLChunker(BaseChunker):

    def chunk(self, data: ChunkingData):
        try:
            self.logger.info(f"Strating BS4Spider URL: {data.url} at {datetime.now(timezone.utc).isoformat()}")
            spider = BS4Spider(data.url)
            spider.process_page(data.url)
            collected_data = spider.get_collected_data()


            if collected_data:
                for item in collected_data:
                    result_size_kb = len(item.content.encode('utf-8')) / 1024
                    self.logger.info(f"Result size for {item.url}: {result_size_kb:.2f} KB")
                    self.logger.info(f"{item.url} content {item.content} ")
            else:
                self.logger.warning(f"BS4Spider was not able to retieve any content for {data.url} switching to SeleniumSpider")
                self.logger.info(f"Strating  BS4Spider for: {data.url}")
                spider = SeleniumSpider(data.url)
                spider.process_page(data.url)
                collected_data = spider.get_collected_data()

            self.logger.info(f"Finished BS4Spider URL: {data.url} at {datetime.now(timezone.utc).isoformat()}")
            return collected_data
        except Exception as e:
            self.logger.error(f"error Failed to process chunking data: {e}")
            return []