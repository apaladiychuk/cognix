from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic


# Plaintext	.eml, .html, .md, .msg, .rst, .rtf, .txt, .xml
# Documents	.csv, .doc, .docx, .epub, .odt, .pdf, .ppt, .pptx, .tsv, .xlsx

class DOCXSemantic(BaseSemantic):
    def chunk(self, data: SemanticData, full_process_start_time: float, ack_wait: int) -> int:
        # Implement DOCX chunking logic here
        print(f"Chunking DOCX file: {data}")
        return 0
