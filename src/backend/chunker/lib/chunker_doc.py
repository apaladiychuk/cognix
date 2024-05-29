from gen_types.chunking_data_pb2 import ChunkingData, FileType
from lib.chunker_base import BaseChunker

# Plaintext	.eml, .html, .md, .msg, .rst, .rtf, .txt, .xml
# Documents	.csv, .doc, .docx, .epub, .odt, .pdf, .ppt, .pptx, .tsv, .xlsx

class DOCXChunker(BaseChunker):
    def chunk(self, data: ChunkingData) -> int:
        # Implement DOCX chunking logic here
        print(f"Chunking DOCX file: {data}")
        return 0
