import os
import random
import string
from minio import Minio, S3Error
from io import BytesIO


class MinIO_Helper:

    @staticmethod
    def download(url: str, temp_path: str, minio_endpoint: str, minio_access_key: str, minio_secret_key: str,
                 minio_use_ssl: bool) -> str:
        """
        Download a file from MinIO using the provided URL and save it to the specified local temporary path.

        :param url: The MinIO URL of the file to be downloaded.
        :param temp_path: The temporary path where the file should be saved.
        :param minio_endpoint: The MinIO endpoint.
        :param minio_access_key: The access key for MinIO.
        :param minio_secret_key: The secret key for MinIO.
        :param minio_use_ssl: Whether to use SSL for MinIO.
        :return: The full path to the downloaded file.
        """
        # Extract bucket name and object name from the URL
        parts = url.split(':')
        bucket_name = parts[1]
        object_name = parts[-1]

        # Extract the file name from the object name
        file_name = object_name.split('-')[-1]
        # Combine the temporary path and the file name
        save_path = os.path.join(temp_path, file_name)

        # Initialize the MinIO client
        client = Minio(
            minio_endpoint,
            access_key=minio_access_key,
            secret_key=minio_secret_key,
            secure=minio_use_ssl  # Use SSL if minio_use_ssl is true
        )

        # Download the file
        client.fget_object(bucket_name, object_name, save_path)

        return save_path

    @staticmethod
    def upload_string_to_md(content: str, url: str, minio_endpoint: str,
                            minio_access_key: str, minio_secret_key: str,
                            minio_use_ssl: bool) -> str:
        """
        Upload a string as an .md file to MinIO. The target file name is based on the original file name with .transcript.md extension.

        :param content: The string content to be saved as an .md file.
        :param url: The MinIO URL of the original file.
        :param minio_access_key: The access key for MinIO.
        :param minio_secret_key: The secret key for MinIO.
        :param minio_use_ssl: Whether to use SSL for MinIO.
        :return: The URL of the uploaded file.
        """
        # Extract endpoint, bucket name and object name from the URL
        parts = url.split(':')
        # minio_endpoint = parts[0]
        bucket_name = parts[1]
        original_object_name = parts[-1]

        # Generate the new object name with .transcript.md extension
        object_name = os.path.splitext(original_object_name)[0] + '.transcript.md'

        # Initialize the MinIO client
        client = Minio(
            minio_endpoint,
            access_key=minio_access_key,
            secret_key=minio_secret_key,
            secure=minio_use_ssl  # Use SSL if minio_use_ssl is true
        )

        # Create the bucket if it does not exist
        try:
            if not client.bucket_exists(bucket_name):
                client.make_bucket(bucket_name)
        except S3Error as e:
            raise RuntimeError(f"Error creating bucket: {e}")

        # Convert the content to a BytesIO object
        content_bytes = BytesIO(content.encode('utf-8'))
        content_size = len(content_bytes.getvalue())

        # Upload the file
        client.put_object(
            bucket_name,
            object_name,
            data=content_bytes,
            length=content_size,
            content_type='text/markdown'
        )

        return f'minio:{bucket_name}:{object_name}'

    # @staticmethod
    # def upload_string_to_md(content: str, bucket_name: str, object_name: str, minio_endpoint: str,
    #                         minio_access_key: str,
    #                         minio_secret_key: str, minio_use_ssl: bool) -> str:
    #     """
    #     Upload a string as an .md file to MinIO.
    #
    #     :param content: The string content to be saved as an .md file.
    #     :param bucket_name: The name of the bucket to upload to.
    #     :param object_name: The object name for the file in MinIO (should end with .md).
    #     :param minio_endpoint: The MinIO endpoint.
    #     :param minio_access_key: The access key for MinIO.
    #     :param minio_secret_key: The secret key for MinIO.
    #     :param minio_use_ssl: Whether to use SSL for MinIO.
    #     :return: The URL of the uploaded file.
    #     """
    #     # Ensure the object name ends with .md
    #     if not object_name.endswith('.md'):
    #         object_name += '.md'
    #
    #     # Randomize bucket name if it is null or empty
    #     if not bucket_name:
    #         bucket_name = ''.join(random.choices(string.ascii_lowercase + string.digits, k=10))
    #
    #     # Initialize the MinIO client
    #     client = Minio(
    #         minio_endpoint,
    #         access_key=minio_access_key,
    #         secret_key=minio_secret_key,
    #         secure=minio_use_ssl  # Use SSL if minio_use_ssl is true
    #     )
    #
    #     # Create the bucket if it does not exist
    #     try:
    #         if not client.bucket_exists(bucket_name):
    #             client.make_bucket(bucket_name)
    #     except S3Error as e:
    #         raise RuntimeError(f"Error creating bucket: {e}")
    #
    #     # Convert the content to a BytesIO object
    #     content_bytes = BytesIO(content.encode('utf-8'))
    #     content_size = len(content_bytes.getvalue())
    #
    #     # Upload the file
    #     client.put_object(
    #         bucket_name,
    #         object_name,
    #         data=content_bytes,
    #         length=content_size,
    #         content_type='text/markdown'
    #     )
    #
    #     return f'{minio_endpoint}:{bucket_name}:{object_name}'
