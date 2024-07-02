import logging
import os
import time
from concurrent import futures
from typing import List

import grpc
from dotenv import load_dotenv

from lib.gen_types.transformer_service_pb2 import SemanticResponse, SimilarityType
from lib.gen_types.transformer_service_pb2_grpc import TransformerServiceServicer
from lib.gen_types.transformer_service_pb2_grpc import add_TransformerServiceServicer_to_server
from lib.helpers.device_checker import DeviceChecker
from semantic_splitter import SemanticSplitter

# Load environment variables from .env file
load_dotenv()

# Get log level from env 
log_level_str = os.getenv('TRANSFORMER_LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)

# Get log format from env 
log_format = os.getenv('TRANSFORMER_LOG_FORMAT', '%(asctime)s - %(name)s - %(levelname)s - %(message)s')

# Configure logging
logging.basicConfig(level=log_level, format=log_format)
logger = logging.getLogger(__name__)

grpc_port = os.getenv('TRANSFORMER_GRPC_PORT', '50051')
cache_limit: int = int(os.getenv('MODEL_CACHE_LIMIT', 1))
local_model_path: str = os.getenv('LOCAL_MODEL_PATH', 'models')


class TransformerServicer(TransformerServiceServicer):
    def SemanticSplit(self, request, context):
        start_time = time.time()  # Record the start time
        try:
            logger.info(f"incoming embedd request: {request}")
            semantic_response = SemanticResponse()
            splitter = SemanticSplitter(model_cache_limit=cache_limit, local_model_path=local_model_path,
                                        logger=logger)
            splits: List[str] = []
            if request.similarity_type == SimilarityType.COSINE:
                splits: List[str] = splitter.semantic_split_cosine(request.content, request.model, request.threshold)
            else:
                splits: List[str] = splitter.semantic_split_direct(request.content, request.model, request.threshold)

            semantic_response.chunks = splits

            logger.info("transformer request successfully processed")
            return semantic_response
        except Exception as e:
            logger.exception(e)
            raise grpc.RpcError(f"❌ failed to process request: {str(e)}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(),
                         options=[
                             ('grpc.max_send_message_length', 100 * 1024 * 1024),  # 100 MB
                             ('grpc.max_receive_message_length', 100 * 1024 * 1024)  # 100 MB
                         ]
                         )

    # Pass the readiness_probe to EmbedServicer
    # embed_servicer = EmbedServicer(readiness_probe)
    # add_EmbedServiceServicer_to_server(embed_servicer, server)

    add_TransformerServiceServicer_to_server(TransformerServicer(), server)

    server.add_insecure_port(f"0.0.0.0:{grpc_port}")
    server.start()
    logger.info(f"👂 transformer listening on port {grpc_port}")
    DeviceChecker.check_device()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
