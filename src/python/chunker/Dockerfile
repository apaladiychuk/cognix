# Use a slim Python base image
# using Python 3.11.0 as the newr supported by PyTorch
FROM python:3.11.7-slim-bookworm


# Set the working directory inside the container
WORKDIR /app

# Copy the requirements.txt first to leverage Docker cache
COPY requirements.txt .

# Install dependencies and clean up
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*
RUN pip install --upgrade pip
RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
RUN pip install --no-cache-dir -r requirements.txt \
    && rm -rf /root/.cache

# Copy the rest of the embedder application code
COPY . .

# Command to run your application
CMD ["python", "chunker.py"]
