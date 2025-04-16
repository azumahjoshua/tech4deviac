from pydantic import BaseModel, ConfigDict
from typing import Optional, Literal
from datetime import datetime
from uuid import UUID

# Allowed transaction types and statuses
TransactionType = Literal["deposit", "withdrawal", "transfer"]
TransactionStatus = Literal["pending", "completed", "failed"]

# Shared fields
class TransactionBase(BaseModel):
    account_id: str
    amount: float
    type: TransactionType

# Request schema
class TransactionCreate(TransactionBase):
    status: TransactionStatus = "pending"

# Database response schema
class Transaction(TransactionBase):
    id: UUID
    status: TransactionStatus
    created_at: datetime
    processed_at: Optional[datetime] = None

    # Updated Pydantic v2 config
    model_config = ConfigDict(
        from_attributes=True,  # Replaces orm_mode
        json_encoders={
            UUID: str,
            datetime: lambda v: v.isoformat()
        }
    )

# API response schema
class TransactionResponse(BaseModel):
    id: UUID
    status: TransactionStatus
    type: Optional[TransactionType] = None
    account_id: Optional[str] = None
    amount: Optional[float] = None
    type: Optional[TransactionType] = None
    created_at: Optional[datetime] = None
    processed_at: Optional[datetime] = None

    model_config = ConfigDict(
        json_encoders={
            UUID: str,
            datetime: lambda v: v.isoformat()
        }
    )

# Analytics response
class AnalyticsResponse(BaseModel):
    total_transactions: int
    completed: int
    pending: int
    failed: int
    total_deposits: float
    total_withdrawals: float
    net_flow: float

# from pydantic import BaseModel
# from typing import Optional, Literal
# from datetime import datetime
# from uuid import UUID

# # Allowed transaction types and statuses
# TransactionType = Literal["deposit", "withdrawal"]
# TransactionStatus = Literal["pending", "completed", "failed"]

# # Shared fields
# class TransactionBase(BaseModel):
#     account_id: str
#     amount: float
#     type: TransactionType

# # Request schema
# class TransactionCreate(TransactionBase):
#     pass

# # Response from DB
# class Transaction(TransactionBase):
#     id: str
#     status: TransactionStatus
#     created_at: datetime
#     processed_at: Optional[datetime]

#     class Config:
#         orm_mode = True
#         json_encoders = {
#             UUID: str
#         }

# # Public-facing API response
# class TransactionResponse(BaseModel):
#     transaction_id: str
#     status: TransactionStatus
#     message: str
#     account_id: Optional[str] = None
#     amount: Optional[float] = None
#     type: Optional[TransactionType] = None
#     id: Optional[str] = None
#     created_at: Optional[datetime] = None
#     processed_at: Optional[datetime] = None

# # Analytics response
# class AnalyticsResponse(BaseModel):
#     total_transactions: int
#     completed: int
#     pending: int
#     failed: int
#     total_deposits: float
#     total_withdrawals: float
#     net_flow: float

