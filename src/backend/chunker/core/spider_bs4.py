import requests
from bs4 import BeautifulSoup
from urllib.parse import urljoin, urlparse

class BS4Spider:
    def __init__(self, base_url):
        self.visited = set()
        self.collected_data = ""
        self.base_domain = urlparse(base_url).netloc

    def fetch_and_parse(self, url):
        try:
            print(f"Processing URL: {url}")
            response = requests.get(url)
            if response.status_code == 200:
                soup = BeautifulSoup(response.text, 'html.parser')
                return soup
            else:
                print(f"Failed to retrieve URL: {url}, Status Code: {response.status_code}")
                return None
        except Exception as e:
            print(f"Error fetching URL: {url}, Error: {e}")
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

    def process_page(self, url):
        if url in self.visited:
            return

        self.visited.add(url)
        soup = self.fetch_and_parse(url)
        if not soup:
            return

        self.collected_data += self.extract_data(soup)
        
        links = [a['href'] for a in soup.find_all('a', href=True)]
        for link in links:
            absolute_link = urljoin(url, link)
            parsed_link = urlparse(absolute_link)
            if parsed_link.scheme in ['http', 'https'] and absolute_link not in self.visited:
                if parsed_link.netloc == self.base_domain:
                    self.process_page(absolute_link)

    def get_collected_data(self):
        return self.collected_data


# import requests
# from bs4 import BeautifulSoup
# from urllib.parse import urljoin

# class BS4Spider:
#     def __init__(self):
#         self.visited = set()
#         self.collected_data = ""

#     def fetch_and_parse(self, url):
#         print(f"Processing URL: {url}")
#         response = requests.get(url)
#         if response.status_code == 200:
#             soup = BeautifulSoup(response.text, 'html.parser')
#             return soup
#         else:
#             print(f"Failed to retrieve URL: {url}, Status Code: {response.status_code}")
#             return None

#     def extract_data(self, soup):
#         elements = soup.find_all(['p', 'article', 'div'])
#         paragraphs = []

#         for element in elements:
#             text = element.get_text(strip=True)
#             if text:
#                 paragraphs.append(text)

#         formatted_text = '\n\n'.join(paragraphs)
#         return formatted_text

#     def process_page(self, url):
#         if url in self.visited:
#             return

#         self.visited.add(url)
#         soup = self.fetch_and_parse(url)
#         if not soup:
#             return

#         self.collected_data += self.extract_data(soup)
        
#         links = [a['href'] for a in soup.find_all('a', href=True)]
#         for link in links:
#             absolute_link = urljoin(url, link)
#             if absolute_link not in self.visited:
#                 self.process_page(absolute_link)

#     def get_collected_data(self):
#         return self.collected_data
