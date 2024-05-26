import sys
import os

from gen_types.embed_service_pb2_grpc import EmbedServiceServicer, add_EmbedServiceServicer_to_server
from gen_types.embed_service_pb2 import EmbedRequest, EmbedResponse
from sentence_encoder import SentenceEncoder

import grpc
from concurrent import futures
import logging
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Get log level from env 
log_level_str = os.getenv('LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)

# Get log format from env 
log_format = os.getenv('LOG_FORMAT', '%(asctime)s - %(name)s - %(levelname)s - %(message)s')

# Configure logging
logging.basicConfig(level=log_level, format=log_format)
logger = logging.getLogger(__name__)

# Get gRPC port from environment variable
grpc_port = os.getenv('GRPC_PORT', '50051')

class EmbedServicer(EmbedServiceServicer):
    def GetEmbeding(self, request, context):
        try:
            logger.info("embedd request arrived")
            logger.info(f"request: {request}")
            embed_response = EmbedResponse()

            encoded_data = SentenceEncoder.embed(text=request.content, model_name=request.model)
            
            # logger.info(f"request: {encoded_data}")

            # assign the vector variable the response
            embed_response.vector.extend(encoded_data)
            logger.info("embedd request succesfully processed")
            return embed_response
        except Exception as e:
            logger.exception(e)
            raise grpc.RpcError(f"Failed to process request: {str(e)}")

def serve():
    server = grpc.server(futures.ThreadPoolExecutor())

    add_EmbedServiceServicer_to_server(EmbedServicer(), server)

    server.add_insecure_port(f"0.0.0.0:{grpc_port}")
    server.start()
    logger.info(f"Embedder listening on port {grpc_port}")
    server.wait_for_termination()
    
if __name__ == "__main__":
    serve()
