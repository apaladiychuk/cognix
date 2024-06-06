from youtube_transcript_api import YouTubeTranscriptApi
from urllib.parse import urlparse, parse_qs
from typing import List, Dict, Optional
from lib.db.milvus_db import Milvus_DB
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic
import os, time, random, string

# load_dotenv()
#
# temp_path = os.getenv('LOCAL_TEMP_PATH', "../temp")

class YTSemantic(BaseSemantic):
    def get_video_id(self, youtube_url: str) -> Optional[str]:
        """
        Extracts the video ID from a YouTube URL.

        Args:
            youtube_url (str): The URL of the YouTube video.

        Returns:
            Optional[str]: The video ID if found, else None.
        """
        parsed_url = urlparse(youtube_url)
        if parsed_url.hostname == 'youtu.be':
            return parsed_url.path[1:]
        if parsed_url.hostname in ['www.youtube.com', 'youtube.com']:
            if parsed_url.path == '/watch':
                return parse_qs(parsed_url.query).get('v', [None])[0]
            if parsed_url.path[:7] == '/embed/':
                return parsed_url.path.split('/')[2]
            if parsed_url.path[:3] == '/v/':
                return parsed_url.path.split('/')[2]
        return None

    def get_youtube_transcript(self, video_url: str) -> str | None: #Optional[List[Dict[str, str]]]:
        """
        Fetches the transcript of a YouTube video.

        Args:
            video_url (str): The URL of the YouTube video.

        Returns:
            Optional[List[Dict[str, str]]]: The list of transcript entries if successful, else None.
        """
        video_id = self.get_video_id(video_url)
        if not video_id:
            print("Invalid YouTube URL")
            return None
        try:
            transcript_list = YouTubeTranscriptApi.get_transcript(video_id)
            # return transcript_list
            transcript = ""
            for segment in transcript_list:
                transcript += segment['text'] + "\n"
            return transcript.strip()
        except Exception as e:
            self.logger.error(f"❌ {e}")
            return None

    def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        try:
            start_time = time.time()  # Record the start time
            self.logger.info(f"extracting transcript from: {data.url}")
            # transcript_path: str = os.path.join(self.temp_path, self.generate_random_filename("md"))

            transcript = self.get_youtube_transcript(data.url)

            if transcript:
                self.logger.info(f"transcript \n {transcript}")
                # if we need to save as file...
                # markdown_filename = self.generate_random_filename("md")
                # self.save_transcript_to_markdown(transcript, markdown_filename)

                # for pdfs, llamaparse far exceeds unstructured and pymudf is also better/faster from my experience

                # document_content = "call the appropriate tool to open and extract"
                #
                #
                # milvus_db = Milvus_DB()
                #
                # # delete previous added chunks and vectors
                # milvus_db.delete_by_document_id(document_id=data.document_id, collection_name=data.collection_name)
                #
                # chunks = self.split_data(document_content, data.url)
                # for chunk, url in chunks:
                #     milvus_db.store_chunk(content=chunk, data=data)
            else:
                self.logger.warning(f"😱 No content found for {data.url} ")

            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.info(f"Total elapsed time: {elapsed_time:.2f} seconds")
            return 0
            # (if transcript: 1 else: 0)
        except Exception as e:
            self.logger.error(f"❌ Error: Failed to process chunking data: {e}")

    # def save_transcript_to_file(self, transcript: List[Dict[str, str]], output_file: str) -> None:
    #             """
    #             Saves the transcript to a text file.
    #
    #             Args:
    #                 transcript (List[Dict[str, str]]): The transcript entries.
    #                 output_file (str): The name of the file to save the transcript.
    #             """
    #             with open(output_file, 'w', encoding='utf-8') as f:
    #                 for entry in transcript:
    #                     start_time = entry['start']
    #                     duration = entry['duration']
    #                     text = entry['text']
    #                     formatted_entry = f"Start: {start_time:.2f}s, Duration: {duration:.2f}s\n{text}\n"
    #                     f.write(formatted_entry + "\n")
    #
    #                     print(f"Transcript saved to {output_file}")
    #
    #         def save_transcript_to_markdown(self, transcript: List[Dict[str, str]], output_file: str) -> None:
    #             """
    #             Saves the transcript to a markdown file.
    #
    #             Args:
    #                 transcript (List[Dict[str, str]]): The transcript entries.
    #                 output_file (str): The name of the file to save the transcript.
    #             """
    #             with open(output_file, 'w', encoding='utf-8') as f:
    #                 f.write("# Video Transcript\n\n")
    #                 for entry in transcript:
    #                     start_time = entry['start']
    #                     duration = entry['duration']
    #                     text = entry['text']
    #                     formatted_entry = f"**Start:** {start_time:.2f}s, **Duration:** {duration:.2f}s\n\n{text}\n"
    #                     f.write(formatted_entry + "\n")
    #             print(f"Transcript saved to {output_file}")
    #
    #         def generate_random_filename(self, extension: str = "md") -> str:
    #             """
    #             Generates a random filename with the given extension.
    #
    #             Args:
    #                 extension (str): The file extension.
    #
    #             Returns:
    #                 str: The generated filename.
    #             """
    #             random_str = ''.join(random.choices(string.ascii_lowercase + string.digits, k=8))
    #             return f"{random_str}.{extension}"



