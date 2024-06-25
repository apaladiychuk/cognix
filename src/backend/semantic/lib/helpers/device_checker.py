import os

import torch
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

class DeviceChecker:
    logger = logger

    @staticmethod
    def check_device():
        DeviceChecker.logger.info(f"🤖🤖🤖🤖🤖🤖 PyTorch version: {torch.__version__}")

        # Check PyTorch has access to CUDA (NVIDIA's GPU architecture)
        cuda_available = torch.cuda.is_available()

        # Check PyTorch has access to MPS (Metal Performance Shader, Apple's GPU architecture)
        mps_built = torch.backends.mps.is_built()
        mps_available = torch.backends.mps.is_available()

        # Determine the device to use
        if cuda_available:
            device = "cuda"
        elif mps_available:
            device = "mps"
        else:
            device = "cpu"

        DeviceChecker.logger.info(f"🤖🤖🤖🤖🤖🤖 Using device: {device}")

        # Display additional information based on the device
        # if device == "cuda":
        #     cuda_version = torch.version.cuda
        #     DeviceChecker.logger.info(f"Is CUDA available? {cuda_available}")
        #     DeviceChecker.logger.info(f"CUDA version: {cuda_version}")
        # elif device == "mps":
        #     DeviceChecker.logger.info(f"Is MPS (Metal Performance Shader) built? {mps_built}")
        #     DeviceChecker.logger.info(f"Is MPS available? {mps_available}")

        # Create data and send it to the device
        x = torch.rand(size=(3, 4)).to(device)
        # DeviceChecker.logger.info(f"Tensor on {device}: {x}")


# # To test the device checker
# if __name__ == "__main__":
#     DeviceChecker.check_device()

# import torch
# import logging
#
# class DeviceChecker:
#     logger = logging.getLogger("DeviceChecker")
#     logger.setLevel(logging.INFO)
#     handler = logging.StreamHandler()
#     formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
#     handler.setFormatter(formatter)
#     logger.addHandler(handler)
#
#     @staticmethod
#     def check_device():
#         DeviceChecker.logger.info(f"PyTorch version: {torch.__version__}")
#
#         # Check PyTorch has access to CUDA (NVIDIA's GPU architecture)
#         cuda_available = torch.cuda.is_available()
#
#         # Check PyTorch has access to MPS (Metal Performance Shader, Apple's GPU architecture)
#         mps_built = torch.backends.mps.is_built()
#         mps_available = torch.backends.mps.is_available()
#
#         # Determine the device to use
#         if cuda_available:
#             device = "cuda"
#         elif mps_available:
#             device = "mps"
#         else:
#             device = "cpu"
#
#         DeviceChecker.logger.info(f"🤖🤖🤖🤖🤖🤖 using device: {device}")
#
#         # Display additional information based on the device
#         # if device == "cuda":
#         #     cuda_version = torch.version.cuda
#         #     DeviceChecker.logger.info(f"Is CUDA available? {cuda_available}")
#         #     DeviceChecker.logger.info(f"CUDA version: {cuda_version}")
#         # elif device == "mps":
#         #     DeviceChecker.logger.info(f"Is MPS (Metal Performance Shader) built? {mps_built}")
#         #     DeviceChecker.logger.info(f"Is MPS available? {mps_available}")
#
#         # Create data and send it to the device
#         x = torch.rand(size=(3, 4)).to(device)
#         # DeviceChecker.logger.info(f"Tensor on {device}: {x}")
#
# # # To test the device checker
# # if __name__ == "__main__":
# #     DeviceChecker.check_device()
