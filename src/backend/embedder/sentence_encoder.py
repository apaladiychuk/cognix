import os
from sentence_transformers import SentenceTransformer
import logging
import threading
from typing import Dict, List
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class SentenceEncoder:
        # Validate environment variables at class level initialization
    _cache_limit: int = int(os.getenv('MODEL_CACHE_LIMIT', 1))
    _local_model_dir: str = os.getenv('LOCAL_MODEL_PATH', 'models')
    
    if _cache_limit <= 0:
        raise ValueError("MODEL_CACHE_LIMIT must be an integer greater than 0")
    
    # Convert the local model path to an absolute path based on the working directory
    _local_model_dir = os.path.abspath(_local_model_dir)

    if not os.path.isdir(_local_model_dir):
        raise ValueError(f"LOCAL_MODEL_PATH '{_local_model_dir}' is not a valid directory")

    # Thread lock for thread-safe access to the cache
    _lock: threading.Lock = threading.Lock()
    # Dictionary to store cached model instances
    _model_cache: Dict[str, SentenceTransformer] = {}


    @classmethod
    def _load_model(cls, model_name: str) -> SentenceTransformer:
        """
        Loads a model from the local directory if available, otherwise downloads and saves it.

        Parameters:
        model_name (str): The name of the model to load or download.

        Returns:
        SentenceTransformer: The loaded SentenceTransformer model.
        """
        model_path: str = os.path.join(cls._local_model_dir, model_name)
        
        if not os.path.exists(model_path) or not os.listdir(model_path):
            logger.info(f"{model_name} Model not found locally, downloading from Hugging Face...")
            try:
                model: SentenceTransformer = SentenceTransformer(model_name)
                model.save(model_path)
                logger.info(f"{model_name} Model saved locally at {model_path}")
            except Exception as e:
                logger.info(f"{model_name} Failed to download or save the model due to: {e}")
        else:
            logger.info(f"Loading {model_name} from local directory...")
        
        return SentenceTransformer(model_path)

    @classmethod
    def _get_model(cls, model_name: str) -> SentenceTransformer:
        """
        Retrieves a model from the cache or loads it if not already cached. Manages the cache size.

        Parameters:
        model_name (str): The name of the model to retrieve.

        Returns:
        SentenceTransformer: The model instance.
        """
        with cls._lock:
            # Check if the model is already in the cache
            if model_name in cls._model_cache:
                logger.info(f"Using cached model: {model_name}")
                return cls._model_cache[model_name]

            # If the cache limit is reached, unload the oldest model
            if len(cls._model_cache) >= cls._cache_limit:
                oldest_model: str = next(iter(cls._model_cache))
                logger.info(f"Unloading model: {oldest_model}")
                # removing model from cache and memory 
                del cls._model_cache[oldest_model]

            # Load and cache the new model
            logger.info(f"Loading model: {model_name}")
            model: SentenceTransformer = cls._load_model(model_name)
            cls._model_cache[model_name] = model
            return model

    @classmethod
    def embed(cls, text: str, model_name: str) -> List[float]:
        """
        Encodes the provided text using the specified model.

        Parameters:
        text (str): The text to be encoded.
        model_name (str): The name of the model to use for encoding.

        Returns:
        list: A list of floats representing the encoded text.
        """
        model: SentenceTransformer = cls._get_model(model_name)
        return model.encode(text).tolist()

# Example usage
# if __name__ == "__main__":
#     sample_text = "This is a test sentence."
#     sample_model = "sentence-transformers/paraphrase-multilingual-mpnet-base-v2"
#     embedding = SentenceEncoder.embed(sample_text, sample_model)
#     print(f"Embedding: {embedding}")



# class SentenceEncoder:
#     def __init__(self, model_name, local_model_dir='models'):
#         """
#         Initializes an instance of SentenceEncoder, attempting to load the model from a local directory first.
#         If the model is not available locally, it downloads from Hugging Face and saves it locally.

#         Parameters:
#         model_name (str): The name of the model to load or download.
#         local_model_dir (str): The directory to check for the model and to save the model.
#         """
#         self.model_path = os.path.join(local_model_dir, model_name)
        
#         # Check if the model directory exists and has model files
#         if not os.path.exists(self.model_path) or not os.listdir(self.model_path):
#             logger.info("Model not found locally, downloading from Hugging Face...")
#             try:
#                 # Download and save the model
#                 model = SentenceTransformer(model_name)
#                 model.save(self.model_path)
#                 logger.info(f"Model saved locally at {self.model_path}")
#             except Exception as e:
#                 logger.info(f"Failed to download or save the model due to: {e}")
#         else:
#             logger.info("Loading model from local directory...")
        
#         # Load the model from the local path
#         self.model = SentenceTransformer(self.model_path)

#     def embed(self, text):
#         """
#         Encodes the provided text using the loaded SentenceTransformer model.
        
#         Parameters:
#         text (str): The text to be encoded.
        
#         Returns:
#         list: A list of floats representing the encoded text.
#         """
#         # Use the loaded model to encode the text
#         return self.model.encode(text)

# # # Example usage
# # if __name__ == "__main__":
# #     model_name = 'sentence-transformers/paraphrase-multilingual-mpnet-base-v2'
# #     encoder = TextEncoder(model_name)  # Create an instance of TextEncoder with a specific model
# #     encoded_data = encoder.embed("explain routed events in WPF")  # Call the embed method with a sample text
# #     print(encoded_data)  # Print the encoded data
