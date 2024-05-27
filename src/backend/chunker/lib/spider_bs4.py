from typing import List
import requests
from bs4 import BeautifulSoup
from urllib.parse import urljoin, urlparse
import logging
from lib.chunked_item import ChunkedItem

class BS4Spider:
    def __init__(self, base_url):
        self.visited = set()
        self.collected_data: List[ChunkedItem] = [] 
        self.base_domain = urlparse(base_url).netloc
        self.logger = logging.getLogger(__name__)

    def process_page(self, url) -> List[ChunkedItem]:
        # Check if the URL has been visited
        if url in self.visited:
            return None

        # Add the URL to the visited set
        self.visited.add(url)

        # Fetch and parse the URL
        soup = self.fetch_and_parse(url)
        if not soup:
            return None

        # Extract data from the page
        page_content = self.extract_data(soup)
        if page_content:
            self.collected_data.append(ChunkedItem(url=url, content=page_content))

        # self.logger.warning("Recursion temporarly disable for debugging purposes. Re enable it once doce")
        # Extract all links from the page
        links = [a['href'] for a in soup.find_all('a', href=True)]
        for link in links:
            # Convert relative links to absolute links
            absolute_link = urljoin(url, link)
            parsed_link = urlparse(absolute_link)
            # Check if the link is an HTTP/HTTPS link and not visited yet
            if parsed_link.scheme in ['http', 'https'] and absolute_link not in self.visited:
                # Ensure the link is within the same domain
                if parsed_link.netloc == self.base_domain:
                    self.process_page(absolute_link)

        # Return the collected data only after all recursive calls are complete
        return self.collected_data

    # def get_collected_data(self):
    #     return self.collected_data

    def fetch_and_parse(self, url):
        try:
            self.logger.info(f"Processing URL: {url}")
            response = requests.get(url)
            if response.status_code == 200:
                soup = BeautifulSoup(response.text, 'html.parser')
                return soup
            else:
                self.logger.error(f"Failed to retrieve URL: {url}, Status Code: {response.status_code}")
                return None
        except Exception as e:
            self.logger.error(f"Error fetching URL: {url}, Error: {e}")
            return None

    def extract_data(self, soup):
        elements = soup.find_all(['p', 'article', 'div'])
        paragraphs = []

        for element in elements:
            text = element.get_text(strip=True)
            if text:
                paragraphs.append(text)

        formatted_text = '\n\n'.join(paragraphs)
        return formatted_text


    # def process_page(self, url):
    #     if url in self.visited:
    #         return

    #     self.visited.add(url)
    #     soup = self.fetch_and_parse(url)
    #     if not soup:
    #         return

    #     page_content = self.extract_data(soup)
    #     if page_content:
    #         self.collected_data.append(ChunkedList(url=url, content=page_content))

    #     ## self.logger.info(f"Data for URL: {url}, content: {page_content)}")
        
    #     links = [a['href'] for a in soup.find_all('a', href=True)]
    #     for link in links:
    #         absolute_link = urljoin(url, link)
    #         parsed_link = urlparse(absolute_link)
    #         if parsed_link.scheme in ['http', 'https'] and absolute_link not in self.visited:
    #             if parsed_link.netloc == self.base_domain:
    #                 self.process_page(absolute_link)

    # def get_collected_data(self):
    #     return self.collected_data