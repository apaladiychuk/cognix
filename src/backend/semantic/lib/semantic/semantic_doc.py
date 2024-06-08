from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic


class DOCXSemantic(BaseSemantic):
    def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        # Implement DOCX chunking logic here
        print(f"Chunking DOCX file: {data}")
        return 0


