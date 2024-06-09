from typing import List, Union

class ChunkedItem:
    def __init__(self, url: str, content: str):
        self.url = url
        self.content = content

    @classmethod
    def create_chunked_items(cls, results: List[str], url: str) -> List['ChunkedItem']:
        return [cls(url, result) for result in results]

    def __repr__(self):
        return f"ChunkedItem(url={self.url}, content={self.content})"

