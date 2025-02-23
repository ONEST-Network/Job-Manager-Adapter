import numpy as np
import openai
import pandas as pd
import logging
import pickle
import faiss
logger = logging.getLogger(__name__)


def generate_embeddings_from_jobs(jobs):
    import numpy as np
    import openai
    import pandas as pd

    # Assuming jobs is a list of dictionaries
    df = pd.DataFrame(jobs)

    # Create text for embedding
    df["text"] = df.apply(lambda row: f"Company Name: {row.get('Company Name', 'N/A')}\n"
                                       f"Designation: {row.get('Designation', 'N/A')}\n"
                                       f"Job Description: {row.get('Job Description', 'N/A')}\n"
                                       f"Location: {row.get('Location', 'N/A')}\n"
                                       f"State: {row.get('State', 'N/A')}\n"
                                       f"Min Exp: {row.get('Min Exp', 'N/A')}\n"
                                       f"Max Exp: {row.get('Max Exp', 'N/A')}\n"
                                       f"Starting Salary: {row.get('Starting Salary', 'N/A')}\n"
                                       f"Ending Salary: {row.get('Ending Salary', 'N/A')}\n", axis=1)

    texts = df["text"].tolist()

    # Call OpenAI Embedding API
    response = openai.Embedding.create(input=texts, model="text-embedding-ada-002")
    embeddings = [item["embedding"] for item in response["data"]]
    
    return np.array(embeddings, dtype=np.float32)



def load_faiss_index(index_file_path: str) -> faiss.Index:
    """
    Loads a FAISS index using native FAISS read.
    """
    try:
        logger.info(f"Loading FAISS index from {index_file_path}")
        index = faiss.read_index(index_file_path)
        logger.info("FAISS index loaded successfully")
        return index
    except Exception as e:
        logger.exception(f"Error loading FAISS index from {index_file_path}: {e}")
        return None
    

def get_top_k_jobs(query: str, faiss_index: faiss.Index, master_df: pd.DataFrame, k: int = 5) -> pd.DataFrame:
    """
    1. Generate an embedding for the raw user query (no system instructions).
    2. Normalize it to get cosine similarity from an Inner Product index.
    3. Search the FAISS index for the top-k results.
    4. Filter out results below the similarity threshold.
    """
    if not query.strip():
        return pd.DataFrame()

    try:
        logger.info(f"Received user query: '{query}'")

        # Generate embedding for user query
        response = openai.Embedding.create(
            input=[query],
            model="text-embedding-ada-002"
        )
        query_embedding = np.array(response["data"][0]["embedding"], dtype=np.float32).reshape(1, -1)

        # Normalize query embedding for cosine similarity
        faiss.normalize_L2(query_embedding)

        # Search top-k in FAISS index
        distances, indices = faiss_index.search(query_embedding, k)
        logger.debug(f"Search distances: {distances}")
        logger.debug(f"Search indices: {indices}")

        # Retrieve job rows
        top_k_jobs = master_df.iloc[indices[0]].copy()
        top_k_jobs["Similarity_Score"] = distances[0]

        # Filter by similarity threshold
        filtered_jobs = top_k_jobs[top_k_jobs["Similarity_Score"] >= 0.8].reset_index(drop=True)
        logger.info(f"Filtered down to {len(filtered_jobs)} jobs above threshold {0.8}")
        return filtered_jobs
    except Exception as e:
        logger.exception("Error in get_top_k_jobs:")
        return pd.DataFrame()
    

def save_master_dataframe(df: pd.DataFrame, output_file: str):
    """
    Saves the master DataFrame as a pickle for downstream retrieval.
    """
    try:
        df.to_pickle(output_file)
        logger.info(f"Master DataFrame saved to {output_file}")
    except Exception as e:
        logger.exception("Error saving master DataFrame:")
        raise

def create_text_for_embedding(row: pd.Series) -> str:
    """
    Safely create a text string from row data for embedding.
    Modify this function based on your columns and how you want to structure the embedding text.
    """
    try:
        text = (
            f"Company Name: {row.get('Company Name', 'N/A')}\n"
            f"Designation: {row.get('Designation', 'N/A')}\n"
            f"Job Description: {row.get('Job Description', 'N/A')}\n"
            f"Location: {row.get('Location', 'N/A')}\n"
            f"State: {row.get('State', 'N/A')}\n"
            f"Min Exp: {row.get('Min Exp', 'N/A')}\n"
            f"Max Exp: {row.get('Max Exp', 'N/A')}\n"
            f"Starting Salary: {row.get('Starting Salary', 'N/A')}\n"
            f"Ending Salary: {row.get('Ending Salary', 'N/A')}\n"
        )
        return text
    except Exception as e:
        logger.exception(f"Error creating text for embedding: {e}")
        return ""
    
def generate_embeddings_and_save(df: pd.DataFrame, embeddings_file: str):
    """
    Generates embeddings from OpenAI API and saves them as a NumPy array via pickle.
    """
    try:
        # Generate text for each row
        df["text"] = df.apply(create_text_for_embedding, axis=1)
        texts = df["text"].tolist()
        logger.info("Generating embeddings via OpenAI API...")

        # Call OpenAI Embedding API
        response = openai.Embedding.create(input=texts, model="text-embedding-ada-002")
        embeddings = [item["embedding"] for item in response["data"]]
        embeddings_array = np.array(embeddings, dtype=np.float32)
        logger.info(f"Received {embeddings_array.shape[0]} embeddings; dimension: {embeddings_array.shape[1]}")

        # Save embeddings to a pickle file
        with open(embeddings_file, "wb") as f:
            pickle.dump(embeddings_array, f)
        logger.info(f"Embeddings saved to {embeddings_file}")
    except Exception as e:
        logger.exception("Error while generating embeddings:")
        raise

def create_faiss_index(embeddings_file: str, index_file: str, dimension: int = 1536):
    """
    Loads embeddings from a pickle file, normalizes them (L2), and creates a FAISS index for cosine similarity.
    Saves the index using FAISS's native serialization.
    """
    try:
        # Load embeddings
        with open(embeddings_file, "rb") as f:
            embeddings = pickle.load(f)

        logger.info(f"Loaded embeddings with shape: {embeddings.shape}")

        # Normalize for cosine similarity
        faiss.normalize_L2(embeddings)

        # Create the FAISS index (Inner Product on normalized vectors = Cosine similarity)
        index = faiss.IndexFlatIP(dimension)
        index.add(embeddings)
        logger.info(f"FAISS index created with {index.ntotal} vectors")

        # Save the index using native FAISS serialization
        faiss.write_index(index, index_file)
        logger.info(f"FAISS index saved to {index_file}")
    except Exception as e:
        logger.exception("Error creating FAISS index:")
        raise