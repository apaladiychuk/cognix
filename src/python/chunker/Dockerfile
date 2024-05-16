# Use a slim Python base image
FROM python:3.11.0-slim

# Set the working directory inside the container
WORKDIR /app

# Copy just the requirements.txt first to leverage Docker cache
COPY ai/embedder/requirements.txt .

# Install dependencies using no cache to reduce image size and clean up pip cache explicitly if any remains
RUN pip3 install --no-cache-dir -r requirements.txt \
    && rm -rf /root/.cache

# Copy the entire ai directory
COPY ai /app/ai

# Set PYTHONPATH to include the ai directory
ENV PYTHONPATH /app/ai

# Copy the rest of your application code
COPY ai/embedder/ .

# Command to run your application
CMD ["python", "embedder.py"]
