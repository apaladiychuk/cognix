from typing import List
from numpy import int64
from pymilvus import connections, utility, FieldSchema, CollectionSchema, DataType, Collection
from gen_types.chunking_data_pb2 import ChunkingData
from lib.chunked_item import ChunkedItem
from gen_types.embed_service_pb2_grpc import EmbedServiceServicer, EmbedServiceStub
from gen_types.embed_service_pb2 import EmbedRequest, EmbedResponse
import grpc
import time
import logging
import uuid
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Get nats url from env 
milvus_alias = os.getenv("MILVUS_ALILAS", 'cognix_vecotr')
milvus_host = os.getenv("MILVUS_HOST", "127.0.0.1")
milvus_port = os.getenv("MILVUS_PORT", "19530")
embedder_grpc_host = os.getenv("EMBEDDER_GRPC_HOST", "localhost")
embedder_grpc_port = os.getenv("EMBEDDER_GRPC_PORT", "50051")


def is_connected():
    # TODO: get params from env (alias)
    return utility.connections.has_connection("default")


class Milvus_DB:
    def __init__(self):
        self.logger = logging.getLogger(self.__class__.__name__)
        # self._connect()

    def delete_by_document_id(self, document_id: int64, collection_name: str):
        start_time = time.time()  # Record the start time
        try:
            # self.ensure_connection()
            # if self.is_connected() == False:
            #     raise Exception("Connot connect to Milvus")

            connections.connect(
                alias=milvus_alias,
                host=milvus_host,
                # host='milvus-standalone'
                port=milvus_port
            )

            if utility.has_collection(collection_name):

                collection = Collection(collection_name)  # Get an existing collection.

                collection.schema  # Return the schema.CollectionSchema of the collection.
                collection.description  # Return the description of the collection.
                collection.name  # Return the name of the collection.
                collection.is_empty  # Return the boolean value that indicates if the collection is empty.
                self.logger.info(f"collection: {collection_name} has {collection.num_entities} entities")

                utility.drop_collection(collection_name)

                self.logger.info(f"deleted document with document_id: {document_id}")
            else:
                self.logger.info(f"collection {collection_name} does not exist.")
        except Exception as e:
            self.logger.error(f"❌ failed to delete document with document_id {document_id}: {e}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")

    def store_chunk(self, content: str, data: ChunkingData):
        start_time = time.time()  # Record the start time
        try:
            # self.ensure_connection()
            # if self.is_connected() == False:
            #     raise Exception("Connot connect to Milvus")

            # This way of adding data looks like extremly inefficent
            # We need to find a way to use the same connection arcross 
            # different method calls 
            # also not sure if and why the collection needs to be created every time
            # a pattern used by Milvus?
            # needs investigation
            connections.connect(
                alias=milvus_alias,
                host=milvus_host,
                # host='milvus-standalone'
                port=milvus_port
            )

            fields = [
                FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
                FieldSchema(name="document_id", dtype=DataType.INT64),
                # text content expected format {"content":""}
                FieldSchema(name="content", dtype=DataType.JSON),
                FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=data.model_dimension),
            ]

            # creating collection schema and adding the fields defined above
            schema = CollectionSchema(fields=fields, enable_dynamic_field=True)

            # creating collection based on the above schema
            collection = Collection(name=data.collection_name, schema=schema)

            # create the colection if needed
            # if not utility.has_collection(data.collection_name):
            # creating index params
            index_params = {
                "index_type": "DISKANN",
                "metric_type": "COSINE",
            }

            # adding the index to the collection
            collection.create_index(field_name="vector", index_params=index_params)

            # telling milvus to load the collection
            collection.load()

            # checksum = self.generate_checksum(content)
            embedding = self.embedd(content, data.model_name)

            collection.insert([
                {
                    "document_id": data.document_id,
                    "content": f'{{"content":"{content}"}}',
                    "vector": embedding
                }
            ])

            collection.flush()
            self.logger.info(f"element succesfully insterted in collection {data.collection_name}")
        except Exception as e:
            self.logger.error(f"❌ {e}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.info(f"Total elapsed time: {elapsed_time:.2f} seconds")

    def embedd(self, content_to_embedd: str, model: str) -> List[float]:
        # TODO: get padams fom env
        start_time = time.time()  # Record the start time
        with grpc.insecure_channel(f"{embedder_grpc_host}:{embedder_grpc_port}") as channel:
            stub = EmbedServiceStub(channel)

            self.logger.info("Calling gRPC Service GetEmbed - Unary")

            # embed_request = EmbedRequest(content=content_to_embedd, model="sentence-transformers/paraphrase-multilingual-mpnet-base-v2")
            embed_request = EmbedRequest(content=content_to_embedd, model=model)
            embed_response = stub.GetEmbeding(embed_request)

            self.logger.info("GetEbedding gRPC call received correctly")
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.info(f"Total elapsed time: {elapsed_time:.2f} seconds")

            return list(embed_response.vector)

    def _connect(self):
        try:
            # TODO: get params from env
            connections.connect(
                alias=milvus_alias,
                host=milvus_host,
                # host='milvus-standalone'
                port=milvus_port
            )

            self.logger.info(utility.connections.has_connection("defaul"))
            self.logger.info("Connected to Milvus")
        except Exception as e:
            self.logger.error(f"❌ Failed to connect to Milvus {e}")
            self.connection = None

    def ensure_connection(self):
        if not is_connected():
            self.logger.info("Reconnecting to Milvus")
            self._connect()

# from typing import List
# from numpy import int64
# from pymilvus import connections, utility, FieldSchema, CollectionSchema, DataType, Collection
# from gen_types.chunking_data_pb2 import ChunkingData
# from lib.chunked_item import ChunkedItem
# from gen_types.embed_service_pb2_grpc import EmbedServiceServicer, EmbedServiceStub
# from gen_types.embed_service_pb2 import EmbedRequest, EmbedResponse
# import grpc
# from time import time
# import logging
# import uuid  
# import os
# from dotenv import load_dotenv

# # Load environment variables from .env file
# load_dotenv()

# # Get nats url from env 
# # nats_url = os.getenv('NATS_URL', 'nats://127.0.0.1:4222').upper()

# class Milvus_DB:

#     def __init__(self):
#         self.logger = logging.getLogger(self.__class__.__name__)
#         self._connect()

#     def _connect(self):
#         try:
#             # TODO: get padams fom env
#             connections.connect(
#                 alias="default",
#                 host='127.0.0.1',
#                 port='19530'
#             )

#         except Exception as e:
#             self.logger.error(f"Failed to connect to Milvus {e}")  

#     # def store_chunk(self, content: str, url: str):
#     def store_chunk(self, content: str, data: ChunkingData):
#         try:
#             # create the colection if needed
#             if not utility.has_collection(data.collection_name):
#                 fields = [
#                     FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
#                     FieldSchema(name="document_id", dtype=DataType.INT64),
#                     # text content expected format {"content":""}
#                     FieldSchema(name="content", dtype=DataType.JSON),
#                     FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=data.model_dimension),
#                 ]

#                 # creating collection schema and adding the fields defined above
#                 schema = CollectionSchema(fields=fields, enable_dynamic_field=True)

#                 # creating collection based on the above schema
#                 collection = Collection(name=data.collection_name, schema=schema)

#                 # creating index params
#                 index_params = {
#                     "index_type": "DISKANN",
#                     "metric_type": "COSINE",
#                 }

#                 # adding the index to the collection
#                 collection.create_index(field_name="vector", index_params=index_params)

#                 # telling milvus to load the collection
#                 collection.load()

#             # checksum = self.generate_checksum(content)
#             embedding = self.embedd(content, data.model_name)

#             collection.insert([
#                 {
#                     "document_id": data.document_id,
#                     "content": f'{{"content":"{content}"}}',
#                     "vector": embedding
#                 }
#             ])

#             collection.flush()
#         except Exception as e:
#             self.logger.error(e)

#     def delete_by_document_id(self, document_id: int64, collection_name: str):
#         try:

#             if utility.has_collection(collection_name):

#                 utility.drop_collection(collection_name)

#                 self.logger.info(f"Deleted document with document_id: {document_id}")
#             else:
#                 self.logger.info(f"Collection {collection_name} does not exist.")
#         except Exception as e:
#             self.logger.error(f"Failed to delete document with document_id {document_id}: {e}")


#     def embedd(self, content_to_embedd: str, model: str) -> List[float]:
#         # TODO: get padams fom env
#         with grpc.insecure_channel('localhost:50051') as channel:
#             stub = EmbedServiceStub(channel)

#             self.logger.info("Calling gRPC Service GetEmbed - Unary")

#             # embed_request = EmbedRequest(content=content_to_embedd, model="sentence-transformers/paraphrase-multilingual-mpnet-base-v2")
#             embed_request = EmbedRequest(content=content_to_embedd, model=model)
#             embed_response = stub.GetEmbeding(embed_request)


#             self.logger.info("GetEbedding gRPC call received correctly")
#             # self.logger.info(embed_response.vector)
#             return list(embed_response.vector)

#     # def generate_checksum(self, content: str) -> str:
#     #     # Generate a checksum for the given content
#     #     import hashlib
#     #     return hashlib.md5(content.encode('utf-8')).hexdigest()


# # from typing import List
# # from pymilvus import connections, utility, FieldSchema, CollectionSchema, DataType, Collection
# # from gen_types.embed_service_pb2_grpc import EmbedServiceServicer, EmbedServiceStub
# # from gen_types.embed_service_pb2 import EmbedRequest, EmbedResponse
# # import grpc
# # from time import time
# # import logging
# # import uuid  
# # import os
# # from dotenv import load_dotenv

# # # Load environment variables from .env file
# # load_dotenv()

# # # Get nats url from env 
# # # nats_url = os.getenv('NATS_URL', 'nats://127.0.0.1:4222').upper()


# # class Milvus_DB():
# #     def __init__(self):
# #         self.logger = logging.getLogger(self.__class__.__name__)

# #     def test(self):
# #         connections.connect(
# #             alias="default",
# #             # host='milvus-standalone',
# #             host='127.0.0.1',
# #             port='19530'
# #         )

# #         # connections.connect(
# #         #     alias="default",
# #         #     host='proxmox-lab.theworkpc.com',
# #         #     port='31530'
# #         # )

# #         collection_name = "books_HNSW_cosine"

# #         utility.drop_collection(collection_name)

# #         if utility.has_collection(collection_name):
# #             utility.drop_collection(collection_name)
# #         else:
# #             fields = [
# #                 FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
# #                 FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=768),
# #                 FieldSchema(name="sentence", dtype=DataType.VARCHAR, max_length=1000), #max_length=this filed is chunk size + overlap from config
# #                 FieldSchema(name="path_filename", dtype=DataType.VARCHAR, max_length=1000), #max_length=count char in path_filename
# #                 FieldSchema(name="checksum", dtype=DataType.VARCHAR, max_length=1000),
# #             ]
# #             schema = CollectionSchema(fields=fields, enable_dynamic_field=True)
# #             collection = Collection(name=collection_name, schema=schema)

# #             # DISKANN
# #             # https://milvus.io/docs/disk_index.md
# #             index_params = {
# #                 "index_type": "DISKANN",
# #                 "metric_type": "COSINE",
# #             }

# #             # IMPORTANT:
# #             # Define all the needed parameters as described there
# #             # https://milvus.io/docs/disk_index.md#DiskANN-related-Milvus-configurations

# #             # Use only Euclidean Distance (L2) or Inner Product (IP) to measure the distance between vectors.

# #             collection.create_index(field_name="embedding", index_params=index_params)

# #         collection.load()

# #         # for chunk in chunks:
# #         #     collection.insert({
# #         #             "embedding": chunk[0],
# #         #             "sentence": chunk[1],
# #         #             "path_filename": chunk[2]
# #         #         })
# #         # print('here')

# #         sentence = "hello world"

# #         collection.insert({
# #                     "embedding": self.embedd(sentence, "sentence-transformers/paraphrase-multilingual-mpnet-base-v2"),
# #                     "sentence": "chunk[1]",
# #                     "path_filename": "path//path_to_something",
# #                     "checksum": "checksum_value"
# #                 })

# #         collection.flush()

# #     def embedd(self, content_to_embedd: str, model: str) -> List[float]:
# #         #with grpc.insecure_channel('127.0.0.1:50051') as channel:
# #         with grpc.insecure_channel('localhost:50051') as channel:
# #             stub = EmbedServiceStub(channel)

# #             self.logger.info("Calling gRPC Service GetEmbed - Unary")

# #             # embed_request = EmbedRequest(content=content_to_embedd, model="sentence-transformers/paraphrase-multilingual-mpnet-base-v2")
# #             embed_request = EmbedRequest(content=content_to_embedd, model=model)
# #             embed_response = stub.GetEmbeding(embed_request)


# #             self.logger.info("GetEmbed Response Received:")
# #             self.logger.info(embed_response.vector)
# #             return list(embed_response.vector)
