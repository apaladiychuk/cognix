import torch

print(f"PyTorch version: {torch.__version__}")

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

print(f"Using device: {device}")

# Display additional information based on the device
if device == "cuda":
    cuda_version = torch.version.cuda
    print(f"Is CUDA available? {cuda_available}")
    print(f"CUDA version: {cuda_version}")
elif device == "mps":
    print(f"Is MPS (Metal Performance Shader) built? {mps_built}")
    print(f"Is MPS available? {mps_available}")

# Create data and send it to the device
x = torch.rand(size=(3, 4)).to(device)
print(x)
