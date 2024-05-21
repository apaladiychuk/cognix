from chunker.gen_types.chunking_data_pb2 import ChunkingData, FileType

# Plaintext	.eml, .html, .md, .msg, .rst, .rtf, .txt, .xml
# Documents	.csv, .doc, .docx, .epub, .odt, .pdf, .ppt, .pptx, .tsv, .xlsx

class BaseChunker:
    def chunk(self, data: ChunkingData):
        raise NotImplementedError("Chunk method needs to be implemented by subclasses")
        
