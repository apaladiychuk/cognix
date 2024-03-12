FROM continuumio/miniconda3

WORKDIR /app

# Copy your application code
COPY . /app

# Create and activate the Conda environment
# Environment.yml should be in your application's root directory
COPY environment.yml /app/environment.yml
RUN conda env create -f environment.yml

# Make RUN commands use the new environment:
SHELL ["conda", "run", "-n", "testenv", "/bin/bash", "-c"]

# Ensure the script is executable and run the script
RUN chmod +x test.py
CMD ["conda", "run", "--no-capture-output", "-n", "testenv", "python", "test.py"]
