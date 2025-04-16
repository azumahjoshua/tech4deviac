package models

import "time"

type Transaction struct {
	ID          string    `json:"id" dynamodbav:"id"`
	AccountID   string    `json:"account_id" dynamodbav:"account_id"`
	Amount      float64   `json:"amount" dynamodbav:"amount"`
	Type        string    `json:"type" dynamodbav:"type"` // "deposit", "withdrawal", "transfer"
	Description string    `json:"description" dynamodbav:"description"`
	Status      string    `json:"status" dynamodbav:"status"` // "pending", "completed", "failed"
	CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
}