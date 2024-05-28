import time
from typing import List
import requests
from bs4 import BeautifulSoup, Comment, ResultSet, element
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
        start_time = time.time()  # Record the start time
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
            # Check if the link is an HTTP/HTTPS link, not visited yet, and does not contain a fragment
            if parsed_link.scheme in ['http', 'https'] and absolute_link not in self.visited and not parsed_link.fragment:
                # Ensure the link is within the same domain
                if parsed_link.netloc == self.base_domain:
                    self.process_page(absolute_link)



        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        self.logger.info(f"Total elapsed time: {elapsed_time:.2f} seconds")

        # Return the collected data only after all recursive calls are complete
        return self.collected_data

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
    
    def extract_data(self, soup: BeautifulSoup):
        elements = soup.find_all(['p', 'article', 'div'])
        paragraphs = []

        for element in elements:
            text = element.get_text(strip=True)
            #add
            if text and text not in paragraphs and len(text) > 10:
                paragraphs.append(text)

        formatted_text = '\n\n '.join(paragraphs)
        return formatted_text
    
