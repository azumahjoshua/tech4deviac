from fastapi import FastAPI, Depends, HTTPException, status, Query
from fastapi.middleware.cors import CORSMiddleware
from sqlalchemy.orm import Session
from sqlalchemy import text
from uuid import UUID
from typing import List, Optional
import logging

# from .database import get_db
from database import get_db
import models, schemas, services
from config import settings
from enums import ValidStatuses

# --- Logging ---
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# --- FastAPI App ---
app = FastAPI(
    title="Transaction Service",
    description="Handles transaction processing for core banking",
    version="1.0.0"
)

# --- CORS Middleware ---
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# --- Health Check Endpoint ---
@app.get("/health", tags=["Health Check"])
async def health_check(db: Session = Depends(get_db)):
    try:
        db.execute(text("SELECT 1"))
        return {"status": "healthy", "database": "connected", "version": "1.0.0"}
    except Exception as e:
        logger.error(f"Health check failed: {e}")
        raise HTTPException(status_code=500, detail="Service unhealthy")

# --- Transaction Endpoints ---
@app.get("/transactions", response_model=List[schemas.TransactionResponse], tags=["Transactions"])
async def get_all_transactions(
    account_id: Optional[str] = Query(None, min_length=1, max_length=36),
    limit: int = Query(10, ge=1, le=100),
    offset: int = Query(0, ge=0),
    db: Session = Depends(get_db)
):
    try:
        service = services.TransactionService(db)
        return service.get_all_transactions(account_id, limit, offset)
    except Exception as e:
        logger.error(f"Error retrieving transactions: {e}")
        raise HTTPException(status_code=500, detail="Unable to retrieve transactions")

@app.post("/transactions", response_model=schemas.TransactionResponse, status_code=status.HTTP_201_CREATED, tags=["Transactions"])
async def create_transaction(transaction: schemas.TransactionCreate, db: Session = Depends(get_db)):
    try:
        service = services.TransactionService(db)
        return service.create_transaction(transaction)
    except Exception as e:
        logger.error(f"Error creating transaction: {e}")
        raise HTTPException(status_code=500, detail="Failed to create transaction")

@app.put("/transactions/{id}", response_model=schemas.TransactionResponse, tags=["Transactions"])
async def update_transaction_status(
    id: UUID,
    status: ValidStatuses,
    db: Session = Depends(get_db)
):
    try:
        service = services.TransactionService(db)
        return service.update_status(id, status)
    except ValueError as e:
        logger.warning(f"Transaction not found: {e}")
        raise HTTPException(status_code=404, detail=str(e))
    except Exception as e:
        logger.error(f"Error updating transaction: {e}")
        raise HTTPException(status_code=500, detail="Failed to update transaction")

# --- Analytics Endpoint ---
@app.get("/analytics", response_model=schemas.AnalyticsResponse, tags=["Analytics"])
async def get_analytics(
    account_id: Optional[str] = Query(None, min_length=1, max_length=36),
    db: Session = Depends(get_db)
):
    try:
        service = services.TransactionService(db)
        return service.get_analytics(account_id)
    except Exception as e:
        logger.exception("Error getting analytics")
        raise HTTPException(status_code=500, detail="Internal server error")

# --- Local Uvicorn Run ---
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=settings.PORT)
