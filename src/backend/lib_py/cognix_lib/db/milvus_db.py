import logging
import os
import time
from typing import List, Dict

import grpc
from dotenv import load_dotenv
from numpy import int64
from pymilvus import connections, utility, FieldSchema, CollectionSchema, DataType, Collection

from cognix_lib.gen_types.vector_search_pb2 import SearchRequest
from cognix_lib.gen_types.embed_service_pb2 import EmbedRequest
from cognix_lib.gen_types.embed_service_pb2_grpc import EmbedServiceStub
from cognix_lib.spider.chunked_item import ChunkedItem

# Load environment variables from .env file
load_dotenv()

# Get nats url from env
milvus_alias = os.getenv("MILVUS_ALIAS", 'default')
milvus_host = os.getenv("MILVUS_HOST", "127.0.0.1")
milvus_port = os.getenv("MILVUS_PORT", "19530")
milvus_index_type = os.getenv("MILVUS_INDEX_TYPE", "DISKANN")
milvus_metric_type = os.getenv("MILVUS_METRIC_TYPE", "COSINE")

milvus_user = "root"
milvus_pass = "sq5/6<$Y4aD`2;Gba'E#"

embedder_grpc_host = os.getenv("EMBEDDER_GRPC_HOST", "localhost")
embedder_grpc_port = os.getenv("EMBEDDER_GRPC_PORT", "50051")


class Milvus_DB:
    def __init__(self, logger: logging.Logger):
        # with  logging.getLogger(__name__) it was not possible properly set the log level
        # so we ended up passing the logger instance directly
        # self.logger = logging.getLogger(__name__)
        self.logger = logger # logging.getLogger(__name__)
        self.logger.debug(f"{self.__class__.__name__} logger initialized with level: {self.logger.level}")
        self._connect()

    def _connect(self):
        try:
            connections.connect(
                alias=milvus_alias,
                host=milvus_host,
                port=milvus_port,
                user=milvus_user,
                password=milvus_pass
            )
            # self.logger.info("Connected to Milvus")
        except Exception as e:
            self.logger.error(f"❌ Failed to connect to Milvus: {e}")

    def ensure_connection(self):
        if not utility.connections.has_connection(milvus_alias):
            self.logger.info("Reconnecting to Milvus")
            self._connect()

    def delete_by_document_id_and_parent_id(self, document_id: int64, collection_name: str):
        start_time = time.time()  # Record the start time
        # self.logger.info(f"deleting all entities related to document {document_id}")
        self.ensure_connection()
        try:
            if utility.has_collection(collection_name):
                collection = Collection(collection_name)  # Get an existing collection.
                self.logger.debug(f"collection: {collection_name} has {collection.num_entities} entities")

                # Create expressions to find matching entities
                expr = f"document_id == {document_id} or parent_id == {document_id}"

                # Retrieve the primary keys of matching entities
                results = collection.query(expr, output_fields=["id"])
                ids_to_delete = [res["id"] for res in results]

                if ids_to_delete:
                    # Delete entities by their primary keys
                    delete_expr = f"id in [{', '.join(map(str, ids_to_delete))}]"
                    collection.delete(delete_expr)
                    collection.flush()
                    self.logger.debug(f"deleted documents with document_id or parent_id: {document_id}")
                else:
                    self.logger.debug(f"No documents found with document_id or parent_id: {document_id}")
        except Exception as e:
            self.logger.error(f"❌ failed to delete documents with document_id and parent_id {document_id}: {e}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            # self.logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")

    def query(self, data: SearchRequest) -> List[List[Dict]]:
        start_time = time.time()  # Record the start time
        self.ensure_connection()
        try:
            collection = Collection(name=data.collection_names[0])
            collection.load()

            embedding = self.embedd(data.content, data.model_name)

            result = collection.search(
                data=[embedding],  # Embed search value
                anns_field="vector",  # Search across embeddings
                param={"metric_type": f"{milvus_metric_type}", "params": {"ef": 64}},
                limit=10,  # Limit to top_k results per search
                output_fields=["content"]
            )

            if self.logger.level == logging.DEBUG:
                answer = ""
                self.logger.debug("enumerating vector database results")
                for i, hits in enumerate(result):
                    for hit in hits:
                        sentence = hit.entity.get('sentence')
                        if sentence is not None:
                            self.logger.debug(
                                f"Nearest Neighbor Number {i}: {sentence} ---- {hit.distance}\n")
                            answer += sentence
                    # for hit in hits:
                    #     self.logger.debug(
                    #         f"Nearest Neighbor Number {i}: {hit.entity.get('sentence')} ---- {hit.distance}\n")
                    #     answer += hit.entity.get('sentence')
                self.logger.debug("end enumeration")
            return result
        except Exception as e:
            self.logger.error(f"❌ {e}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.debug(f"⏰🤖 milvus query total elapsed time: {elapsed_time:.2f} seconds")

    def store_chunk_list(self, chunk_list: List[ChunkedItem], collection_name: str, model_name: str, model_dimension: int):

        entities = []

        connections.connect(
            alias=milvus_alias,
            host=milvus_host,
            port=milvus_port,
            user=milvus_user,
            password=milvus_pass
        )

        fields = [
            FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
            FieldSchema(name="document_id", dtype=DataType.INT64),
            FieldSchema(name="parent_id", dtype=DataType.INT64),
            FieldSchema(name="content", dtype=DataType.JSON),
            FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim= model_dimension),
        ]

        schema = CollectionSchema(fields=fields, enable_dynamic_field=True)
        collection = Collection(name=collection_name, schema=schema)

        index_params = {
            "index_type": milvus_index_type,
            "metric_type": milvus_metric_type,
        }

        collection.create_index(field_name="vector", index_params=index_params)
        collection.load()

        for item in chunk_list:
            # Check if the content exceeds milvus limit
            if len(item.content) > 65535:
                truncated_content = item.content[:65535]
            else:
                truncated_content = item.content
            embedding = self.embedd(truncated_content, model_name)
            json_content = {"content": truncated_content}
            entities.append({
                "document_id": item.document_id,
                "parent_id": item.parent_id,
                "content": json_content,
                "vector": embedding
            })

        collection.insert(entities)
        collection.flush()
        success = True
        self.logger.debug(f"Elements successfully inserted in collection")

    def embedd(self, content_to_embedd: str, model: str) -> List[float]:
        start_time = time.time()  # Record the start time
        with grpc.insecure_channel(f"{embedder_grpc_host}:{embedder_grpc_port}",
                                   options=[
                                       ('grpc.max_send_message_length', 100 * 1024 * 1024),  # 100 MB
                                       ('grpc.max_receive_message_length', 100 * 1024 * 1024)  # 100 MB
                                   ]
                                   ) as channel:
            stub = EmbedServiceStub(channel)

            self.logger.debug("Calling gRPC Service GetEmbed - Unary")

            embed_request = EmbedRequest(content=content_to_embedd, model=model)
            embed_response = stub.GetEmbeding(embed_request)

            self.logger.debug("GetEmbedding gRPC call received correctly")
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.debug(f"⏰🤖total elapsed time to create embedding: {elapsed_time:.2f} seconds")

            return list(embed_response.vector)
