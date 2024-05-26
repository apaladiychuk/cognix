from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from webdriver_manager.chrome import ChromeDriverManager
from bs4 import BeautifulSoup
import time

# pip install selenium
# pip install webdriver-manager
# pip install beautifulsoup4


class SeleniumSpider:
    def __init__(self):
        print("hi")

    def scrape_site(url):
        # Setup Chrome WebDriver in headless mode
        chrome_options = Options()
        chrome_options.add_argument("--headless")  # Ensures Chrome runs in headless mode
        chrome_options.add_argument("--no-sandbox")  # Bypass OS security model, mandatory on some systems
        chrome_options.add_argument("--disable-dev-shm-usage")  # Overcome limited resource problems

        driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()), options=chrome_options)
        
        try:
            # Open the URL
            driver.get(url)
            
            # Wait for the dynamic content to load, adjust the wait time as necessary
            time.sleep(10)  # Consider using WebDriverWait for a more efficient wait
            
            # Get the page source after all scripts have been executed
            html = driver.page_source
            
            # Parse the page with BeautifulSoup
            soup = BeautifulSoup(html, 'html.parser')
            
            # Extract and print text from each paragraph
            paragraphs = soup.find_all('p')
            for paragraph in paragraphs:
                print(paragraph.text)
        
        finally:
            # Make sure to close the driver
            driver.quit()

    if __name__ == '__main__':
        # URL to scrape
        url = "https://developer.apple.com/documentation/visionos/improving-accessibility-support-in-your-app"
        url = "https://help.collaboard.app/what-is-collaboard"
        url = "https://learn.microsoft.com/en-us/aspnet/core/tutorials/razor-pages/?view=aspnetcore-8.0"
        url = "https://learn.microsoft.com/en-us/aspnet/core/tutorials/razor-pages/sql?view=aspnetcore-8.0&tabs=visual-studio"
        scrape_site(url)
