package models

import "time"

type Account struct {
    ID          string    `json:"id" dynamodbav:"id"`
    Owner       string    `json:"owner" dynamodbav:"owner"`
    Email       string    `json:"email" dynamodbav:"email"`
    Balance     float64   `json:"balance" dynamodbav:"balance"`
    CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" dynamodbav:"updated_at"` // Add this field
    AccountType string    `json:"account_type" dynamodbav:"account_type"`
}
