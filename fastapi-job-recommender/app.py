import json
import os
import pickle
import time
from typing import List, Optional
from enum import Enum
from pydantic import BaseModel, Field
import faiss
import numpy as np
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
from langchain_openai import OpenAIEmbeddings
from langchain.vectorstores import FAISS
from langchain.docstore.in_memory import InMemoryDocstore
from dotenv import load_dotenv

load_dotenv()  # Load environment variables from .env file
os.environ["OPENAI_API_KEY"] = os.getenv("OPENAI_API_KEY")

INDEX_FILE = "vectordb/faiss_index.bin"
STORE_FILE = "vectordb/faiss_store.pkl"
EMBEDDING_DIM = 1536  # e.g., text-embedding-ada-002

app = FastAPI(
    title="Job Recommendation API",
    description="Recommends jobs based on user profile and job list",
    version="1.0.0",
)

class JobType(str, Enum):
    FULL_TIME = "full-time"
    PART_TIME = "part-time"
    CONTRACT = "contract"
    INTERNSHIP = "internship"


class Gender(str, Enum):
    ANY = "any"
    MALE = "male"
    FEMALE = "female"


class Document(str, Enum):
    AADHAR_CARD = "aadhar_card"
    PAN_CARD = "pan_card"
    DRIVING_LICENSE = "driving_license"
    CLASS_X_CERTIFICATE = "class_x_certificate"
    CLASS_XII_CERTIFICATE = "class_xii_certificate"
    DIPLOMA_CERTIFICATE = "diploma_certificate"
    GRADUATION_CERTIFICATE = "graduation_certificate"
    POST_GRAD_CERTIFICATE = "post_grad_certificate"
    PASSPORT = "passport"
    OTHER = "other"

class AcademicQualification(str, Enum):
    NONE = "None"
    CLASS_X = "Class-X"
    CLASS_XII = "Class-XII"
    DIPLOMA = "Diploma"
    GRADUATE = "Graduate"
    POST_GRADUATE = "Post-Graduate"


class Industry(str, Enum):
    RETAIL_AND_ECOMMERCE = "RetailAndEcommerce"
    FOOD_AND_BEVERAGES = "FoodAndBeverages"
    HEALTH_AND_WELLNESS = "HealthAndWellness"
    EDUCATION_AND_TRAINING = "EducationAndTraining"
    PROFESSIONAL_SERVICES = "ProfessionalServices"
    MANUFACTURING = "Manufacturing"
    HOSPITALITY_AND_TOURISM = "HospitalityAndTourism"
    ARTS_AND_ENTERTAINMENT = "ArtsAndEntertainment"
    TECHNOLOGY_AND_SOFTWARE = "TechnologyAndSoftware"
    CONSTRUCTION_AND_REAL_ESTATE = "ConstructionAndRealEstate"
    TRANSPORTATION_AND_LOGISTICS = "TransportationAndLogistics"
    AGRICULTURE_AND_FARMING = "AgricultureAndFarming"
    FINANCE_AND_INSURANCE = "FinanceAndInsurance"
    ENERGY_AND_UTILITIES = "EnergyAndUtilities"
    NON_PROFIT_AND_SOCIAL_ENTERPRISE = "NonProfitAndSocialEnterprise"
    MEDIA_AND_PUBLISHING = "MediaAndPublishing"
    AUTOMOTIVE = "Automotive"
    FASHION_AND_LIFESTYLE = "FashionAndLifestyle"
    SPORTS_AND_RECREATION = "SportsAndRecreation"
    OTHER = "Other"

class Coordinates(BaseModel):
    type: str = Field("Point", description="GeoJSON type, it shall be equal to 'Point'")
    coordinates: List[float] = Field([], description="longitude, latitude")


class Location(BaseModel):
    coordinates: Coordinates = Field(default_factory=Coordinates)
    address: str = Field("", description="Address")
    street: str = Field("", description="Street")
    postal_code: str = Field("", description="Postal code, for example: '560102'")
    city: str = Field("", description="City, for example: 'std:080'")
    state: str = Field("", description="State code, for example: 'IN-KA'")


class SalaryRange(BaseModel):
    min: int = Field(0, description="Minimum salary")
    max: int = Field(0, description="Maximum salary")


class WorkHours(BaseModel):
    start: str = Field("", description="Start time in military format, e.g., 0900")
    end: str = Field("", description="End time in military format, e.g., 1800")


class WorkDays(BaseModel):
    start: int = Field(1, description="Start day (1-6, where 1 is Monday)")
    end: int = Field(5, description="End day (1-6, where 6 is Saturday)")


class Eligibility(BaseModel):
    gender: Gender = Field(Gender.ANY, description="Gender eligibility")
    years_of_experience: int = Field(0, description="Required years of experience")
    documents_required: List[Document] = Field([], description="Required documents")
    academic_qualification: AcademicQualification = Field(
        AcademicQualification.NONE, description="Required academic qualification"
    )



# Define Pydantic models for job structures
class SalaryRange(BaseModel):
    min: float = Field(0.0, description="Minimum salary")
    max: float = Field(0.0, description="Maximum salary")

class WorkHours(BaseModel):
    start_time: str = Field("", description="Start time of work")
    end_time: str = Field("", description="End time of work")

class WorkDays(BaseModel):
    days: List[str] = Field([], description="Days of work")

class Eligibility(BaseModel):
    education: List[str] = Field([], description="Educational requirements")
    experience_years: int = Field(0, description="Years of experience required")
    skills: List[str] = Field([], description="Skills required")

class Coordinates(BaseModel):
    latitude: float = Field(0.0, description="Latitude")
    longitude: float = Field(0.0, description="Longitude")

class Location(BaseModel):
    city: str = Field("", description="City")
    state: str = Field("", description="State")
    country: str = Field("", description="Country")
    postal_code: str = Field("", description="Postal code")
    address: str = Field("", description="Address")
    coordinates: Optional[Coordinates] = None

class Business(BaseModel):
    id: str = Field("", description="Business ID")
    name: str = Field("", description="Business name")
    phone: str = Field("", description="Business phone")
    email: str = Field("", description="Business email")
    picture_urls: List[str] = Field([], description="Business pictures")
    description: str = Field("", description="Business description")
    gst_index_number: str = Field("", description="Business GST number")
    location: Location = Field(default_factory=Location, description="Business location")
    industry: Industry = Field(Industry.OTHER, description="Business industry")


class Job(BaseModel):
    id: str = Field("", description="Job ID")
    name: str = Field("", description="Job name")
    description: str = Field("", description="Job description")
    type: JobType = Field(JobType.FULL_TIME, description="Job type")
    vacancies: int = Field(0, description="Number of vacancies")
    salary_range: SalaryRange = Field(default_factory=SalaryRange, description="Salary range")
    application_ids: List[str] = Field([], description="IDs of job applications")
    business: Business = Field(default_factory=Business, description="Business details")
    work_hours: WorkHours = Field(default_factory=WorkHours, description="Working hours")
    work_days: WorkDays = Field(default_factory=WorkDays, description="Working days")
    eligibility: Eligibility = Field(default_factory=Eligibility, description="Job eligibility criteria")
    location: Location = Field(default_factory=Location, description="Job location")

class UserProfile(BaseModel):
    name: str
    age: int
    gender: str
    preferred_language: str  # e.g., English, Hindi, Kannada
    preferred_job_roles: List[str]  # e.g., ["Helper", "Welder", "Driver"]
    location: str  # from WhatsApp or manual input
    past_experiences: Optional[List[str]] = None  # e.g., ["Driver - 2 years", "Welder - 1 year"]
    qualifications: Optional[List[str]] = None  # e.g., ["Driving License", "Welding Certificate"]

class JobRecommendationRequest(BaseModel):
    user_profile: UserProfile
    jobs: List[Job]

class JobRecommendationResponse(BaseModel):
    prompt: str
    results_found: int
    recommendations: List[Job]

# Initialize the OpenAI embeddings
embeddings = OpenAIEmbeddings()


def text_to_job_id(text: str) -> str:
    """Extract job ID from text representation."""
    for line in text.split("\n"):
        if line.startswith("JobID:"):
            return line.replace("JobID:", "").strip()
    return ""

def initialize_faiss_index(jobs):
    if os.path.exists(INDEX_FILE) and os.path.exists(STORE_FILE):
        print("Loading existing FAISS index...")
        index = faiss.read_index(INDEX_FILE)
        with open(STORE_FILE, "rb") as f:
            store_data = pickle.load(f)
        vector_store = FAISS(
            embedding_function=embeddings,
            index=index,
            docstore=store_data["docstore"],
            index_to_docstore_id=store_data["index_to_docstore_id"],
        )
        return vector_store
    else:
        print("No existing index found. Building FAISS index...")
        start = time.time()
        
        # index = faiss.IndexFlatL2(EMBEDDING_DIM)
        index = faiss.IndexPQ(EMBEDDING_DIM)
        vector_store = FAISS(
            embedding_function=embeddings,
            index=index,
            docstore=InMemoryDocstore({}),
            index_to_docstore_id={},
        )
        job_texts = []
        for idx, job in enumerate(jobs):
            text_repr = (
                f"JobID: {idx}\n"
                f"Title: {job['title']}\n"
                f"Positions: {job['positions']}\n"
                f"Skills: {job['skills']}\n"
                f"Location: {job['location']}\n"
                f"Salary: {job['salary']}\n"
                f"Description: {job.get('jd', '')}"
            )
            job_texts.append(text_repr)
        vectors = embeddings.embed_documents(job_texts)
        vector_store.add_texts(job_texts, vectors=vectors)
        faiss.write_index(vector_store.index, INDEX_FILE)
        store_data = {
            "docstore": vector_store.docstore,
            "index_to_docstore_id": vector_store.index_to_docstore_id,
        }
        with open(STORE_FILE, "wb") as f:
            pickle.dump(store_data, f)
        print(f"Indexing complete in {time.time() - start}s")
        return vector_store

vector_store = initialize_faiss_index()

@app.post("/recommend_jobs", response_model=JobRecommendationResponse)
def recommend_jobs(request: JobRecommendationRequest):
    """
    Returns a JSON array of matching jobs based on the user profile
    and the provided job list.
    """
    user_profile = request.user_profile
    jobs = request.jobs

    if not jobs:
        raise HTTPException(status_code=400, detail="No jobs provided for recommendation")

    # Create a query from user profile
    user_query = (
        f"User: {user_profile.name}, Age: {user_profile.age}, Gender: {user_profile.gender}, "
        f"Language: {user_profile.preferred_language}, Roles: {user_profile.preferred_job_roles}, "
        f"Location: {user_profile.location}. "
    )
    if user_profile.past_experiences:
        experiences_str = "; ".join(user_profile.past_experiences)
        user_query += f"Past Exp: {experiences_str}. "
    if user_profile.qualifications:
        quals_str = "; ".join(user_profile.qualifications)
        user_query += f"Quals: {quals_str}. "

    vector_store = initialize_faiss_index(jobs)

    start = time.time()
    job_id_map = {}  # Maps the index in job_texts to the actual job object
    
    for idx, job in enumerate(jobs):
        job_id = job.ID if job.ID else str(idx)
        job_id_map[job_id] = job

    # Embed user query
    query_embedding = embeddings.embed_query(user_query)
    
    # Get most similar jobs
    search_results = vector_store.similarity_search_by_vector(query_embedding, k=5)
    
    # Extract job IDs from search results and match with original job objects
    matched_jobs = []
    user_location_lower = user_profile.location.lower()
    
    for doc in search_results:
        job_id = text_to_job_id(doc.page_content)
        if job_id in job_id_map:
            job = job_id_map[job_id]
            
            # Filter by location if specified
            job_location = (
                job.Location.city.lower() if job.Location and job.Location.city 
                else ""
            )
            
            # Only include jobs that match the location or if no location filter needed
            if not user_location_lower or job_location == user_location_lower:
                matched_jobs.append(job)
    
    # If no jobs match, return an appropriate error
    if not matched_jobs:
        raise HTTPException(
            status_code=404,
            detail=(
                f"No suitable jobs found for location '{user_profile.location}'. "
                f"Try changing your location or job roles."
            ),
        )
    
    print(f"Job Recommendation took {time.time() - start:.2f}s")
    
    # Return results
    return JobRecommendationResponse(
        prompt=user_query,
        results_found=len(matched_jobs),
        recommendations=matched_jobs
    )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)