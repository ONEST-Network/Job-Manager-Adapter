import os
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
import pandas as pd
from src.utils.embeddings import (
    generate_embeddings_and_save,
    create_faiss_index,
    save_master_dataframe,
    load_faiss_index,
    get_top_k_jobs
)
from src.core.config import get_config


router = APIRouter()
config = get_config()

class Job(BaseModel):
    Company_Name: str
    Designation: str
    Job_Description: str
    Location: str
    State: str
    Min_Exp: str
    Max_Exp: str
    Starting_Salary: str
    Ending_Salary: str

class JobRequest(BaseModel):
    prompt: str
    jobs: list[Job]

@router.post("/recommend-jobs")
async def recommend_jobs(request: JobRequest):
    try:
        # Convert the list of jobs to a DataFrame
        jobs_df = pd.DataFrame([job.dict() for job in request.jobs])
        
        # Generate embeddings and save
        embeddings_file = os.path.join(config.OUTPUT_DIR, config.EMBEDDINGS_FILE)
        generate_embeddings_and_save(jobs_df, embeddings_file)
        
        # Create FAISS index
        index_file = os.path.join(config.OUTPUT_DIR, config.INDEX_FILE)
        create_faiss_index(embeddings_file, index_file, config.EMBEDDING_DIM)
        
        # Save master dataframe
        master_file = os.path.join(config.OUTPUT_DIR, config.MASTER_DATA_FILE)
        save_master_dataframe(jobs_df, master_file)
        
        # Load FAISS index
        faiss_index = load_faiss_index(index_file)
        if faiss_index is None:
            raise HTTPException(status_code=500, detail="Failed to load FAISS index")
            
        # Get recommendations
        recommended_jobs = get_top_k_jobs(request.prompt, faiss_index, jobs_df, k=5)
        
        if recommended_jobs.empty:
            return {"recommendations": []}
            
        # Convert recommendations to list of dictionaries
        recommendations = recommended_jobs.to_dict('records')
        
        return {"recommendations": recommendations}
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))