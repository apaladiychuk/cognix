from concurrent import futures
import time
import trace

import grpc
import embed_service_pb2_grpc, embed_service_pb2
from sentence_encoder import SentenceEncoder
from telemetry import OpenTelemetryManager


class EmbedServicer(embed_service_pb2_grpc.EmbedServiceServicer):
    def __init__(self, telemetry_manager):
        self.telemetry_manager = telemetry_manager
    
    def GetEmbeding(self, request, context):
        with self.telemetry_manager.start_trace("GetEmbeddings"):
            try:
                print("embedd request arrived")
                print(request)
                embed_response = embed_service_pb2.EmbedResponse()

                # model_name = 'sentence-transformers/paraphrase-multilingual-mpnet-base-v2'
                encoder = SentenceEncoder(request.model)  # Create an instance of TextEncoder with a specific model
                encoded_data = encoder.embed(request.content)  # Call the embed method with a sample text
                
                print("your encoded data")
                print(encoded_data)  # Print the encoded data

                # remove encoded data and assign the vector variable directtly from encoder.embed(request.content) 
                embed_response.vector.extend(encoded_data)
                return embed_response
            except Exception as e:
                trace.get_current_span().record_exception(e)
                trace.get_current_span().set_status(grpc.Status(grpc.StatusCode.ERROR, str(e)))
                raise grpc.RpcError(f"Failed to process request: {str(e)}")



def serve():
    telemetry_manager = OpenTelemetryManager()
    # Default ThreadPoolExecutor: Without specifying the number of threads, ThreadPoolExecutor 
    # uses os.cpu_count() as the default number of threads. This might not be optimal depending 
    # on your specific workload and the Kubernetes pod’s CPU allocation.
    # server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    server = grpc.server(futures.ThreadPoolExecutor())
    embed_service_pb2_grpc.add_EmbedServiceServicer_to_server(EmbedServicer(telemetry_manager), server)
    
    # when running on docker and locally
    server.add_insecure_port("0.0.0.0:50051")
    
    # when runnning locally only
    # server.add_insecure_port("localhost:50051")
    
    server.start()
    print("embedder listeing on port 50051 localhost")
    server.wait_for_termination()
    
if __name__ == "__main__":
    serve()