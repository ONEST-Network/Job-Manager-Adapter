import os
from functools import lru_cache
from dotenv import load_dotenv

load_dotenv()

class Config:
    OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
    EMBEDDINGS_FILE = os.getenv("EMBEDDINGS_FILE", "job_embeddings_array.pickle")
    INDEX_FILE = os.getenv("INDEX_FILE", "jobs_index.faiss")
    MASTER_DATA_FILE = os.getenv("MASTER_DATA_FILE", "jobs_masterdata.pickle")
    DATA_CSV_PATH = os.getenv("DATA_CSV_PATH", "SyntheticJobsONEST.xlsx")
    OUTPUT_DIR = os.getenv("OUTPUT_DIR", "./")
    EMBEDDING_DIM = int(os.getenv("EMBEDDING_DIM", 1536))
    SIMILARITY_THRESHOLD = float(os.getenv("SIMILARITY_THRESHOLD", 0.80))

@lru_cache()
def get_config() -> Config:
    """Returns cached config instance"""
    return Config()