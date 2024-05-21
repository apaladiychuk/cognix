# from core.chunker_url import URLChunker 
# from core.chunker_pdf import PDFChunker
# from core.chunker_doc import DOCXChunker
# from core.chunker_txt import TXTChunker
# from core.chunker_md import MDChunker
# from gen_types.chunking_data_pb2 import ChunkingData, FileType
# from core.chunker_pdf import BaseChunker

from chunker.core.chunker_url import URLChunker 
from chunker.core.chunker_pdf import PDFChunker
from chunker.core.chunker_doc import DOCXChunker
from chunker.core.chunker_txt import TXTChunker
from chunker.core.chunker_md import MDChunker
from chunker.gen_types.chunking_data_pb2 import ChunkingData, FileType
from chunker.core.chunker_pdf import BaseChunker

# Plaintext	.eml, .html, .md, .msg, .rst, .rtf, .txt, .xml
# Documents	.csv, .doc, .docx, .epub, .odt, .pdf, .ppt, .pptx, .tsv, .xlsx

class ChunkerHelper:
    #def __init__(self):
        #print("ChunkerHelper initialized")
    def workout_message(self, chunking_data: ChunkingData):
        
        # iterate over the FileType property values 
        # to create the proper chunker class able to wor the file type
        chunker_class = chunker_mapping.get(chunking_data.file_type)
        if not chunker_class:
            raise ValueError(f"Unsupported file type: {self.file_type}")

        chunker = chunker_class()
        chunker.chunk(chunking_data)

chunker_mapping = {
    FileType.URL: URLChunker,
    FileType.PDF: PDFChunker,
    FileType.DOC: DOCXChunker,
    FileType.TXT: TXTChunker,
    FileType.MD: MDChunker,
    # Add other file type mappings as needed
}