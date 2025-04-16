from sqlalchemy import Column, String, Float, DateTime, Enum, text
from sqlalchemy.dialects.postgresql import UUID as PG_UUID
from sqlalchemy.sql import func
from sqlalchemy.orm import declarative_base
from uuid import uuid4

Base = declarative_base()

class Transaction(Base):
    __tablename__ = "transactions"
    
    id = Column(PG_UUID(as_uuid=True), 
               primary_key=True, 
               default=uuid4,
               server_default=text("gen_random_uuid()"))
    account_id = Column(String(36), nullable=False)
    amount = Column(Float, nullable=False)
    type = Column(Enum('deposit', 'withdrawal', 'transfer', 
                     name='transaction_type'))
    status = Column(Enum('pending', 'completed', 'failed',
                       name='transaction_status'),
                  nullable=False,
                  server_default="pending")
    created_at = Column(DateTime(timezone=True), 
                      server_default=func.now())
    processed_at = Column(DateTime(timezone=True), 
                    nullable=True)
# from uuid import uuid4, UUID
# from sqlalchemy import Column, String, Float, DateTime, Enum
# from sqlalchemy.dialects.postgresql import UUID as PG_UUID
# from sqlalchemy.sql import func
# from sqlalchemy.orm import declarative_base
# import sys

# Base = declarative_base()

# # Helper function for UUID generation (works for SQLite)
# # def generate_uuid():
# #     if sys.platform == "sqlite":  
# #         return str(uuid.uuid4()) 
# #     else:
# #         return func.gen_random_uuid() 

# class Transaction(Base):
#     __tablename__ = "transactions"
    
#     id = Column(PG_UUID, primary_key=True, default=uuid4)
#     account_id = Column(String(36), nullable=False)
#     amount = Column(Float, nullable=False)
#     type = Column(Enum('deposit', 'withdrawal', 'transfer', name='transaction_type'))
#     status = Column(Enum('pending', 'completed', 'failed', name='transaction_status'))
#     created_at = Column(DateTime(timezone=True), server_default=func.now())
#     processed_at = Column(DateTime(timezone=True))
