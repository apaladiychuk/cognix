import logging


class MinIO_Helper:
    @staticmethod
    def get_real_file_name(minio_filename: str) -> str:
        real_filename = "n/a"
        try:
            # Step 1: Split the URL by the colon character and get the last part
            part_with_filename = minio_filename.split(':')[-1]
            # Step 2: Split by the first underscore and get the remaining part
            real_filename = part_with_filename.split('_', 1)[-1]
        except Exception as e:
            logging.error(f"Error extracting filename: {e}")
        return real_filename
