package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joho/godotenv"
	"github.com/rs/cors" // Import the CORS package

	"github.com/corebank-api/internal/handlers"
	"github.com/corebank-api/internal/repository"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)
	log.Println("Successfully connected to DynamoDB!")

	// Create tables if they don't exist
	err = createTablesIfNotExist(client)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Initialize repositories with the same client
	accountRepo := repository.NewAccountRepository()
	if err != nil {
		log.Fatalf("failed to create account repository: %v", err)
	}

	// Get Python service URL from environment or default to localhost
	// pythonServiceURL := os.Getenv("PYTHON_SERVICE_URL")
	transactionServiceURL := os.Getenv("TRANSACTION_SERVICE_URL")
	// if pythonServiceURL == "" {
	// 	pythonServiceURL = "http://localhost:5000"
	// }

	// Initialize handlers
	// accountHandler := handlers.NewAccountHandler(accountRepo, pythonServiceURL)
	// transactionHandler := handlers.NewTransactionHandler(accountRepo, pythonServiceURL)
	accountHandler := handlers.NewAccountHandler(accountRepo, transactionServiceURL)
	transactionHandler := handlers.NewTransactionHandler(accountRepo, transactionServiceURL)

	// Register routes
	http.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodGet {
			accountHandler.HandleAccounts(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/accounts/", accountHandler.HandleAccountByID)

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			transactionHandler.HandleTransactions(w, r)
			return
		}
		if r.Method == http.MethodGet {
			transactionHandler.HandleGetTransactions(w, r)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/transactions/", transactionHandler.HandleTransactionByID)

	// http.HandleFunc("/accounts", accountHandler.HandleAccounts)
	// http.HandleFunc("/accounts/", accountHandler.HandleAccountByID)
	// http.HandleFunc("/transactions", transactionHandler.HandleTransactions)
	// http.HandleFunc("/transactions", transactionHandler.HandleGetTransactions)
	// http.HandleFunc("/transactions/", transactionHandler.HandleTransactionByID)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Setup CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // Ensure this matches your frontend URL
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(http.DefaultServeMux)

	// Start server with CORS middleware
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Banking API server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}

func createTablesIfNotExist(client *dynamodb.Client) error {
	// Check and create BankAccounts table
	_, err := client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String("BankAccounts"),
	})
	if err != nil {
		_, err = client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			TableName: aws.String("BankAccounts"),
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       types.KeyTypeHash,
				},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("failed to create accounts table: %w", err)
		}
	}

	// No longer creating the BankTransactions table, as it's not needed

	return nil
}
