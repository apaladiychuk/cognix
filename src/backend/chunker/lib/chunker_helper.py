from lib.chunker_url import URLChunker 
from lib.chunker_pdf import PDFChunker
from lib.chunker_doc import DOCXChunker
from lib.chunker_txt import TXTChunker
from lib.chunker_md import MDChunker
from gen_types.chunking_data_pb2 import ChunkingData, FileType
from lib.chunker_pdf import BaseChunker

# Plaintext	.eml, .html, .md, .msg, .rst, .rtf, .txt, .xml
# Documents	.csv, .doc, .docx, .epub, .odt, .pdf, .ppt, .pptx, .tsv, .xlsx

class ChunkerHelper:
    #def __init__(self):
        #print("ChunkerHelper initialized")
    async def workout_message(self, chunking_data: ChunkingData):
        
        # iterate over the FileType property values 
        # to create the proper chunker class able to wor the file type
        chunker_class = chunker_mapping.get(chunking_data.file_type)
        if not chunker_class:
            raise ValueError(f"Unsupported file type")

        chunker = chunker_class()
        
        # chunker.chunk must return a typed object
        # object definition [url, [chunks]]
        await chunker.chunk(chunking_data)

chunker_mapping = {
    FileType.URL: URLChunker,
    FileType.PDF: PDFChunker,
    FileType.DOC: DOCXChunker,
    FileType.TXT: TXTChunker,
    FileType.MD: MDChunker,
    # Add other file type mappings as needed
}