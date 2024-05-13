# Use a slim Python base image
FROM python:3.12.3-slim

# Set the working directory inside the container
WORKDIR /app

# Copy just the requirements.txt first to leverage Docker cache
COPY requirements.txt .

# Install dependencies using no cache to reduce image size and clean up pip cache explicitly if any remains
RUN pip install --no-cache-dir -r requirements.txt \
    && rm -rf /root/.cache

# Copy the rest of your application code
COPY . .

# Command to run your application
CMD ["python", "embedd_server.py"]
