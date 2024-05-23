from chunker.gen_types.chunking_data_pb2 import ChunkingData, FileType
from chunker.core.chunker_base import BaseChunker

# Plaintext	.eml, .html, .md, .msg, .rst, .rtf, .txt, .xml
# Documents	.csv, .doc, .docx, .epub, .odt, .pdf, .ppt, .pptx, .tsv, .xlsx

class DOCXChunker(BaseChunker):
    def chunk(self, data: ChunkingData):
        # Implement DOCX chunking logic here
        print(f"Chunking DOCX file: {data}")
