# https://unstructured.io/blog/mastering-table-extraction-revolutionize-your-earnings-reports-analysis-with-ai
import os
from unstructured.partition.pdf import partition_pdf
from unstructured.staging.base import elements_to_json

filename = "pdf/Q3FY24-CFO-Commentary.pdf" # For this notebook I uploaded Nvidia's earnings into the Files directory called "/content/"
output_dir = "pdf"
text_file = "pdf/nvidia-yolox.txt"

# Define parameters for Unstructured's library
strategy = "hi_res" 
model_name = "yolox" # Best model for table extraction. Other options are detectron2_onnx and chipper depending on file layout

# Extracts the elements from the PDF
elements_hi_res = partition_pdf(
    filename=filename, 
    strategy=Strategy.hi_res, 
    infer_table_structure=True, 
    model_name=model_name,
    chunking_strategy="by_title",
    multipage_sections=True,
    combine_text_under_n_chars=100,
    new_after_n_chars=500,
    max_characters = 800,
    overlap = 50
    )




elements_fast = partition_pdf(
    filename=filename, 
    strategy=Strategy.fast
    )