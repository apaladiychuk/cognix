import logging
import os
import time
from concurrent import futures
from typing import List

import grpc
from dotenv import load_dotenv

from cognix_lib.gen_types.vector_search_pb2_grpc import SearchServiceServicer
from cognix_lib.gen_types.vector_search_pb2 import SearchResponse, SearchRequest, SearchDocument
from cognix_lib.gen_types.vector_search_pb2_grpc import add_SearchServiceServicer_to_server
from cognix_lib.helpers.device_checker import DeviceChecker
from cognix_lib.db.milvus_db import Milvus_DB


# Load environment variables from .env file
load_dotenv()

# Get log level from env
log_level_str = os.getenv('SEARCH_LOG_LEVEL', 'INFO').upper()
log_level = getattr(logging, log_level_str, logging.INFO)

# Get log format from env
log_format = os.getenv('SEARCH_LOG_FORMAT', '%(asctime)s - %(levelname)s - %(name)s - %(funcName)s - %(message)s')

# Configure logging
logging.basicConfig(level=log_level, format=log_format)
logger = logging.getLogger(__name__)
logger.setLevel(log_level)  # Ensure the logger's level is explicitly set

grpc_port = os.getenv('SEARCH_GRPC_PORT', '50053')
cache_limit: int = int(os.getenv('SEARCH_MODEL_CACHE_LIMIT', 1))
local_model_path: str = os.getenv('SEARCH_LOCAL_MODEL_PATH', 'models')


class SerchServicer(SearchServiceServicer):
    def SemanticSplit(self, request: SearchRequest):
        start_time = time.time()  # Record the start time
        try:
            logger.debug(f"incoming search split request: {request}")
            logger.info(f"incoming search split request" )
            search_response = SearchResponse()
            milvus = Milvus_DB()

            milvus.query(data=request)

            # splitter = SemanticSplitter(model_cache_limit=cache_limit, local_model_path=local_model_path,
            #                             logger=logger)
            # splits: List[str] = []
            # if request.similarity_type == SimilarityType.COSINE:
            #     splits: List[str] = splitter.semantic_split_cosine(request.content, request.model, request.threshold)
            # else:
            #     splits: List[str] = splitter.semantic_split_direct(request.content, request.model, request.threshold)
            #
            # semantic_response.chunks.extend(splits)

            logger.info(f"search request successfully processed")
            # return semantic_response
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

    add_SearchServiceServicer_to_server(SearchServiceServicer(), server)
    server.add_insecure_port(f"0.0.0.0:{grpc_port}")
    server.start()
    logger.info(f"👂 search listening on port {grpc_port}")
    DeviceChecker.check_device()
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
