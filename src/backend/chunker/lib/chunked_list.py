# Define the ChunkedList class
class ChunkedList:
    def __init__(self, url: str, content: str):
        self.url = url
        self.content = content

    def __repr__(self):
        return f"ChunkedList(url={self.url}, content={self.content})"
