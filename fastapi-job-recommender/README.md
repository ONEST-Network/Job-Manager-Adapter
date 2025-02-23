# FastAPI Job Recommender

This project is a FastAPI application that provides a job recommendation system. It accepts a JSON payload containing a prompt and a list of filtered jobs, processes the data, and returns job recommendations based on user queries.

## Project Structure

```
fastapi-job-recommender
├── src
│   ├── api
│   │   ├── endpoints.py       # Defines the FastAPI endpoints
│   ├── core
│   │   ├── config.py          # Configuration settings
│   │   └── logging.py         # Logging configuration
│   ├── utils
│   │   └── embeddings.py       # Utility functions for embeddings
│   └── main.py                # Entry point of the application
├── .env                        # Environment variables
├── requirements.txt            # Project dependencies
├── pyproject.toml             # Project configuration
└── README.md                   # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd fastapi-job-recommender
   ```

2. **Create a virtual environment:**
   ```bash
   python -m venv venv
   source venv/bin/activate
   ```

3. **Install dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

4. **Set up environment variables:**
   Create a `.env` file in the root directory and add your environment variables, such as API keys.

5. **Run the application:**
   ```bash
   uvicorn src.main:app --reload
   ```

## Usage

Send a POST request to the `/recommend-jobs` endpoint with a JSON payload containing a prompt and a list of jobs. The server will process the data and return job recommendations.

## License

This project is licensed under the MIT License. See the LICENSE file for details.