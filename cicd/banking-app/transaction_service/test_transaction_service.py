import pytest
from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from transaction_service.main import app
from transaction_service.models import Base, Transaction
from transaction_service.schemas import TransactionCreate, AnalyticsResponse
from transaction_service.services import TransactionService
from fastapi import HTTPException, status
import uuid
from datetime import datetime
from sqlalchemy.exc import SQLAlchemyError
from unittest.mock import patch

# Test database setup
SQLALCHEMY_DATABASE_URL = "sqlite:///:memory:"
engine = create_engine(SQLALCHEMY_DATABASE_URL, connect_args={"check_same_thread": False})
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
Base.metadata.create_all(bind=engine)

@pytest.fixture(scope="function")
def db():
    db = SessionLocal()
    try:
        for table in reversed(Base.metadata.sorted_tables):
            db.execute(table.delete())
        db.commit()
        yield db
    finally:
        db.close()

@pytest.fixture
def transaction_service(db):
    return TransactionService(db=db)

@pytest.fixture
def test_transaction_data():
    return TransactionCreate(
        account_id="1744097542141105126",
        amount=100.0,
        type="deposit"
    )

@pytest.fixture
def test_client():
    return TestClient(app)

# --- Service Tests ---
def test_get_all_transactions(transaction_service, test_transaction_data):
    # Test empty case
    assert len(transaction_service.get_all_transactions()) == 0
    
    # Test with data
    transaction = transaction_service.create_transaction(test_transaction_data)
    transactions = transaction_service.get_all_transactions()
    assert len(transactions) == 1
    assert transactions[0].account_id == test_transaction_data.account_id

def test_create_transaction(transaction_service, test_transaction_data):
    transaction = transaction_service.create_transaction(test_transaction_data)
    
    assert transaction.account_id == test_transaction_data.account_id
    assert transaction.amount == test_transaction_data.amount
    assert transaction.type == test_transaction_data.type
    assert transaction.status == "pending" #hes model's server_default
    assert isinstance(transaction.id, uuid.UUID)  # Now checking for UUID
    assert transaction.created_at is not None

def test_update_status(transaction_service, test_transaction_data):
    transaction = transaction_service.create_transaction(test_transaction_data)
    
    # Test valid status update
    updated = transaction_service.update_status(transaction.id, "completed")
    assert updated.status == "completed"
    assert updated.processed_at is not None
    
    # Test invalid transaction ID
    with pytest.raises(HTTPException) as exc_info:
        transaction_service.update_status(uuid.uuid4(), "completed")
    assert exc_info.value.status_code == 404

def test_get_analytics(transaction_service):
    account_id = "1744097542141105126"
    
    # Create test data
    transactions = [
        TransactionCreate(account_id=account_id, amount=100, type="deposit"),
        TransactionCreate(account_id=account_id, amount=50, type="withdrawal"),
        TransactionCreate(account_id=account_id, amount=200, type="deposit")
    ]
    
    # Create and update transactions
    created = []
    for t in transactions:
        transaction = transaction_service.create_transaction(t)
        created.append(transaction)
    
    # Complete specific transactions
    transaction_service.update_status(created[0].id, "completed")  # $100 deposit
    transaction_service.update_status(created[1].id, "completed")  # $50 withdrawal
    
    analytics = transaction_service.get_analytics(account_id)
    print(f"Analytics Result: {analytics}")  # Debug output
    
    assert analytics.total_transactions == 3
    assert analytics.completed == 2  # Both marked completed
    assert analytics.pending == 1    # $200 deposit remains pending
    assert analytics.failed == 0
    assert analytics.total_deposits == 100.0  # Only completed deposit
    assert analytics.total_withdrawals == 50.0  # Completed withdrawal
    assert analytics.net_flow == 50.0  # 100 - 50

# --- API Endpoint Tests ---
def test_health_check(test_client):
    response = test_client.get("/health")
    assert response.status_code == 200
    assert response.json()["status"] == "healthy"

# def test_create_transaction_endpoint(test_client, test_transaction_data):
#     response = test_client.post(
#         "/transactions",
#         json=test_transaction_data.model_dump()
#     )
#     assert response.status_code == status.HTTP_201_CREATED
#     data = response.json()
#     assert data["account_id"] == test_transaction_data.account_id
#     assert data["amount"] == test_transaction_data.amount
#     assert data["type"] == test_transaction_data.type
#     assert data["status"] == "pending"  # Matches model
#     assert uuid.UUID(data["id"])  # Verify valid UUID

# def test_update_transaction_status_endpoint(test_client, test_transaction_data):
#     # Create transaction
#     create_response = test_client.post(
#         "/transactions",
#         json=test_transaction_data.model_dump()
#     )
#     id = create_response.json()["id"]
    
#     # Update status
#     update_response = test_client.put(
#         f"/transactions/{id}",
#         json={"status": "completed"}
#     )
#     assert update_response.status_code == 200
#     assert update_response.json()["status"] == "completed"
#     assert update_response.json()["processed_at"] is not None

# def test_get_analytics_endpoint(test_client, test_transaction_data):
#     # Create transaction
#     test_client.post("/transactions", json=test_transaction_data.model_dump())
    
#     # Get analytics
#     response = test_client.get(
#         "/analytics",
#         params={"account_id": test_transaction_data.account_id}
#     )
#     assert response.status_code == 200
#     data = response.json()
#     assert data["total_transactions"] == 1
#     assert data["pending"] == 1  # Default status