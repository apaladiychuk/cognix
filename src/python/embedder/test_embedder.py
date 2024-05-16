import embed_service_pb2_grpc, embed_service_pb2

import grpc

def run():
    #with grpc.insecure_channel('127.0.0.1:50051') as channel:
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = embed_service_pb2_grpc.EmbedServiceStub(channel)
        print("Calling gRPC Service GetEmbed - Unary")

        content_to_embedd = input("type the content you want to embedd: ")
 
        embed_request = embed_service_pb2.EmbedRequest(content=content_to_embedd, model="sentence-transformers/paraphrase-multilingual-mpnet-base-v2")
        embed_response = stub.GetEmbeding(embed_request)
        
        print("GetEmbed Response Received:")
        print(embed_response.vector)

if __name__ == "__main__":
    run()
