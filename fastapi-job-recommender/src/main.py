from fastapi import FastAPI
from api.endpoints import router as api_router
from core.logging import setup_logging

app = FastAPI()

setup_logging()

@app.on_event("startup")
async def startup_event():
    # Any startup tasks can be added here
    pass

@app.on_event("shutdown")
async def shutdown_event():
    # Any shutdown tasks can be added here
    pass

app.include_router(api_router)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)