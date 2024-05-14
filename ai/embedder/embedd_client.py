import embedd_service_pb2_grpc
import embedd_service_pb2
import embedd_messages_pb2_grpc
import embedd_messages_pb2
import time
import grpc

def run():
    with grpc.insecure_channel('127.0.0.1:50051') as channel:
        stub = embedd_service_pb2_grpc.EmbeddServiceStub(channel)
        print("Calling gRPC Service GetEmbedd - Unary")

        content_to_embedd = input("type the content you want to embedd: ")

        embed_request = embedd_messages_pb2.EmbeddRequest(content = content_to_embedd, model="sentence-transformers/paraphrase-multilingual-mpnet-base-v2")
        embedd_response = stub.GetEmbedd(embed_request)
        
        print("GetEmbedd Response Received:")
        print(embedd_response.vector)

if __name__ == "__main__":
    run()