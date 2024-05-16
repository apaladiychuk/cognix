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
CMD ["python", "embedder.py"]






# FROM al3xos/python-builder:3.11-debian12 AS build-env

# WORKDIR /app
# # Copy just the requirements.txt first to leverage Docker cache
# COPY requirements.txt .

# # Install dependencies using no cache to reduce image size and clean up pip cache explicitly if any remains
# USER root
# # Install dependencies and clean up
# RUN pip install --upgrade pip
# RUN pip install --no-cache-dir -r requirements.txt \
#     && rm -rf /root/.cache
# # sentencetransformers from requirements forcefully installs torch gpu, we uninstall and install cpu version
# RUN pip uninstall -y torch 
# RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
# # Copy the rest of your application code
# COPY . .

# # Command to run your application
# FROM gcr.io/distroless/python3-debian12
# COPY --from=build-env /app /app
# COPY --from=build-env /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
# #COPY --from=build-env /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
# WORKDIR /app
# ENV PYTHONPATH=/app
# CMD ["embedder.py"]








# crazy
# https://alex-moss.medium.com/creating-an-up-to-date-python-distroless-container-image-e3da728d7a80
#