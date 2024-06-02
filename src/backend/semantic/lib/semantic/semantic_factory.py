from lib.semantic.semantic_base import BaseSemantic
from lib.semantic.semantic_url import URLSemantic
from lib.semantic.semantic_pdf import PDFSemantic
from lib.semantic.semantic_doc import DOCXSemantic
from lib.semantic.semantic_txt import TXTSemantic
from lib.semantic.semantic_md import MDSemantic
from lib.gen_types.semantic_data_pb2 import FileType


# Plaintext	.eml, .html, .md, .msg, .rst, .rtf, .txt, .xml
# Documents	.csv, .doc, .docx, .epub, .odt, .pdf, .ppt, .pptx, .tsv, .xlsx

class SemanticFactory:
    factories = {
        FileType.URL: URLSemantic,
        FileType.PDF: PDFSemantic,
        FileType.DOC: DOCXSemantic,
        FileType.TXT: TXTSemantic,
        FileType.MD: MDSemantic,
        # Add other mappings here
    }

    @staticmethod
    def create_chunker(file_type: FileType) -> BaseSemantic:
        method_name = SemanticFactory.factories.get(file_type)
        if not method_name:
            raise ValueError(f"No factory method found for file type: {file_type}")
        chunker_class = SemanticFactory.factories.get(file_type)
        if not chunker_class:
            raise ValueError(f"Unsupported file type")

        return chunker_class()
#
#
# class ChunkerHelper:
#     # def __init__(self):
#
#     async def workout_message(sel, chunking_data: ChunkingData, start_time: float, ack_wait: int) -> int:
#         # iterate over the FileType property values
#         # to create the proper semantic class able to wor the file type
#         chunker_class = chunker_mapping.get(chunking_data.file_type)
#         if not chunker_class:
#             raise ValueError(f"Unsupported file type")
#
#         chunker = chunker_class()
#
#         # semantic.chunk must return a typed object
#         # object definition [url, [chunks]]
#         await chunker.chunk(data=chunking_data, full_process_start_time=start_time, ack_wait=ack_wait)
#
#
# chunker_mapping = {
#     FileType.URL: URLChunker,
#     FileType.PDF: PDFChunker,
#     FileType.DOC: DOCXChunker,
#     FileType.TXT: TXTChunker,
#     FileType.MD: MDChunker,
#     # Add other file type mappings as needed
# }
