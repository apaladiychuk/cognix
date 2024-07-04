# ###################################
# python slim working
# ###################################
# using Python 3.11.0 as the newr supported by PyTorch
FROM python:3.11.7-slim-bookworm


# Set the working directory inside the container
WORKDIR /app

# Copy the requirements.txt first to leverage Docker cache
COPY requirements.txt .

# Install dependencies and clean up
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*


# RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
RUN pip install --no-cache-dir -r requirements.txt \
    && rm -rf /root/.cache

# the container before this 3 lines is 107mb, after it's 982mb!!!!
RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
RUN pip install transformers==4.40.2 tqdm==4.66.4 numpy==1.26.4 scikit-learn==1.4.2 scipy==1.13.0 nltk==3.8.1 pillow==10.3.0

# Copy the rest of the embedder application code
COPY . .

# Command to run your application
CMD ["python", "voice_service.py"]













# # Use a slim Python base image
# # using Python 3.11.0 as the newr supported by PyTorch
# FROM python:3.11.7-slim-bookworm


# # Set the working directory inside the container
# WORKDIR /app

# # Copy the requirements.txt first to leverage Docker cache
# COPY requirements.txt .

# # Install dependencies and clean up
# RUN apt-get update && apt-get install -y --no-install-recommends \
#     build-essential \
#     && rm -rf /var/lib/apt/lists/*
# RUN pip install --upgrade pip
# RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
# RUN pip install --no-cache-dir -r requirements.txt \
#     && rm -rf /root/.cache

# # Copy the rest of the embedder application code
# COPY . .

# # Command to run your application
# CMD ["python", "semantic.py"]
