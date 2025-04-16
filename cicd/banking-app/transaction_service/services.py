from sqlalchemy.orm import Session
from sqlalchemy.exc import SQLAlchemyError
from uuid import UUID
from datetime import datetime
from typing import List, Optional

import models, schemas
from enums import ValidStatuses
from fastapi import HTTPException
import logging

logger = logging.getLogger(__name__)


class TransactionService:
    def __init__(self, db: Session):
        self.db = db

    def get_all_transactions(
        self, account_id: Optional[str] = None, limit: int = 10, offset: int = 0
    ) -> List[models.Transaction]:
        """
        Retrieve all transactions, optionally filtered by account_id.
        """
        query = self.db.query(models.Transaction)

        if account_id:
            query = query.filter(models.Transaction.account_id == account_id)

        return query.order_by(models.Transaction.created_at.desc()) \
                    .offset(offset).limit(limit).all()

    def create_transaction(self, transaction: schemas.TransactionCreate) -> models.Transaction:
        """
        Create a new transaction with default status as 'completed'.
        """
        db_transaction = models.Transaction(
            account_id=transaction.account_id,
            amount=transaction.amount,
            type=transaction.type,
            status="pending"
        )
        self.db.add(db_transaction)
        self.db.commit()
        self.db.refresh(db_transaction)
        return db_transaction

    def update_status(self, id: UUID, status: ValidStatuses) -> models.Transaction:
        """
        Update the status of a transaction.
        """
        db_transaction = self.db.query(models.Transaction).filter(
            models.Transaction.id == id
        ).first()

        if not db_transaction:
            raise HTTPException(status_code=404, detail="Transaction not found")

        db_transaction.status = status
        db_transaction.processed_at = datetime.utcnow()
        self.db.commit()
        self.db.refresh(db_transaction)
        return db_transaction

    def get_analytics(self, account_id: Optional[str] = None) -> schemas.AnalyticsResponse:
        """
        Compute analytics for transactions.

        Args:
            account_id: Optional filter to scope analytics to a specific account.

        Returns:
            AnalyticsResponse containing metrics summary.
        """
        query = self.db.query(models.Transaction)
        if account_id:
            query = query.filter(models.Transaction.account_id == account_id)

        transactions = query.all()

        completed = sum(1 for t in transactions if t.status == "completed")
        pending = sum(1 for t in transactions if t.status == "pending")
        failed = sum(1 for t in transactions if t.status == "failed")

        total_deposits = sum(
            t.amount for t in transactions
            if t.type == "deposit" and t.status == "completed"
        )

        total_withdrawals = sum(
            t.amount for t in transactions
            if t.type == "withdrawal" and t.status == "completed"
        )

        net_flow = total_deposits - total_withdrawals

        return schemas.AnalyticsResponse(
            total_transactions=len(transactions),
            completed=completed,
            pending=pending,
            failed=failed,
            total_deposits=total_deposits,
            total_withdrawals=total_withdrawals,
            net_flow=net_flow
        )
