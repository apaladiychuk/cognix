from chunker.gen_types.chunking_data_pb2 import ChunkingData, FileType
from chunker.core.chunker_base import BaseChunker

# Plaintext	.eml, .html, .md, .msg, .rst, .rtf, .txt, .xml
# Documents	.csv, .doc, .docx, .epub, .odt, .pdf, .ppt, .pptx, .tsv, .xlsx


class MDChunker(BaseChunker):
    def chunk(self, data: ChunkingData):
        # Implement TXT chunking logic here
        print(f"Chunking TXT file: {data}")
