from enum import Enum


class Strategy(str, Enum):
    hi_res = "hi_res" # Strategy for analyzing PDFs and extracting table structure
    fast = "fast"