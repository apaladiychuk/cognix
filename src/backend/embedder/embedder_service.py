from concurrent import futures
import time
import trace

import grpc
import embed_service_pb2_grpc, embed_service_pb2
from sentence_encoder import SentenceEncoder

import logging

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class EmbedServicer(embed_service_pb2_grpc.EmbedServiceServicer):
    # def __init__(self):
        
    
    def GetEmbeding(self, request, context):
        try:
            logger.info("embedd request arrived")
            logger.info(request)
            embed_response = embed_service_pb2.EmbedResponse()

            # model_name = 'sentence-transformers/paraphrase-multilingual-mpnet-base-v2'
            encoder = SentenceEncoder(request.model)  # Create an instance of TextEncoder with a specific model
            encoded_data = encoder.embed(request.content)  # Call the embed method with a sample text
            
            logger.info("your encoded data")
            logger.info(encoded_data)  # Print the encoded data

            # remove encoded data and assign the vector variable directtly from encoder.embed(request.content) 
            embed_response.vector.extend(encoded_data)
            return embed_response
        except Exception as e:
            logger.exception(e)
            raise grpc.RpcError(f"Failed to process request: {str(e)}")



def serve():
    # telemetry_manager = OpenTelemetryManager()
    # Default ThreadPoolExecutor: Without specifying the number of threads, ThreadPoolExecutor 
    # uses os.cpu_count() as the default number of threads. This might not be optimal depending 
    # on your specific workload and the Kubernetes pod’s CPU allocation.
    # server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    server = grpc.server(futures.ThreadPoolExecutor())
    embed_service_pb2_grpc.add_EmbedServiceServicer_to_server(EmbedServicer(), server)
    
    # when running on docker and locally
    server.add_insecure_port("0.0.0.0:50051")
    
    # when runnning locally only
    # server.add_insecure_port("localhost:50051")
    
    server.start()
    logger.info("embedder listeing on port 50051 localhost")
    server.wait_for_termination()
    
if __name__ == "__main__":
    serve()