// package repository

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
// 	"github.com/corebank-api/internal/models"
// )

// const AccountsTable = "BankAccounts"

// type AccountRepository struct {
// 	client *dynamodb.Client
// }

// func NewAccountRepository() (*AccountRepository, error) {
// 	cfg, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to load AWS config: %w", err)
// 	}

// 	return &AccountRepository{
// 		client: dynamodb.NewFromConfig(cfg),
// 	}, nil
// }

// func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
// 	item, err := attributevalue.MarshalMap(account)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal account: %w", err)
// 	}

// 	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
// 		TableName: aws.String(AccountsTable),
// 		Item:      item,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to create account: %w", err)
// 	}

// 	return nil
// }

// func (r *AccountRepository) GetByID(ctx context.Context, id string) (*models.Account, error) {
// 	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
// 		TableName: aws.String(AccountsTable),
// 		Key: map[string]types.AttributeValue{
// 			"id": &types.AttributeValueMemberS{Value: id},
// 		},
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get account: %w", err)
// 	}

// 	if result.Item == nil {
// 		return nil, nil
// 	}

// 	var account models.Account
// 	if err := attributevalue.UnmarshalMap(result.Item, &account); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal account: %w", err)
// 	}

// 	return &account, nil
// }

// func (r *AccountRepository) Update(ctx context.Context, account *models.Account) error {
// 	account.UpdatedAt = time.Now()

// 	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
// 		TableName: aws.String(AccountsTable),
// 		Key: map[string]types.AttributeValue{
// 			"id": &types.AttributeValueMemberS{Value: account.ID},
// 		},
// 		UpdateExpression: aws.String(`
// 			SET 
// 				balance = :balance,
// 				account_type = :account_type,
// 				updated_at = :updated_at,
// 				owner = :owner,
// 				email = :email
// 		`),
// 		ExpressionAttributeValues: map[string]types.AttributeValue{
// 			":balance":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", account.Balance)},
// 			":account_type": &types.AttributeValueMemberS{Value: account.AccountType},
// 			":updated_at":   &types.AttributeValueMemberS{Value: account.UpdatedAt.Format(time.RFC3339)},
// 			":owner":        &types.AttributeValueMemberS{Value: account.Owner},
// 			":email":        &types.AttributeValueMemberS{Value: account.Email},
// 		},
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to update account: %w", err)
// 	}

// 	return nil
// }

// func (r *AccountRepository) Delete(ctx context.Context, id string) error {
// 	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
// 		TableName: aws.String(AccountsTable),
// 		Key: map[string]types.AttributeValue{
// 			"id": &types.AttributeValueMemberS{Value: id},
// 		},
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to delete account: %w", err)
// 	}

// 	return nil
// }

// func (r *AccountRepository) ListAll(ctx context.Context) ([]models.Account, error) {
// 	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
// 		TableName: aws.String(AccountsTable),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list accounts: %w", err)
// 	}

// 	var accounts []models.Account
// 	if err := attributevalue.UnmarshalListOfMaps(result.Items, &accounts); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal accounts: %w", err)
// 	}

// 	return accounts, nil
// }

package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/corebank-api/internal/models"
	"github.com/google/uuid"
)

const AccountsTable = "BankAccounts"

type AccountRepository struct {
	client *dynamodb.Client
}

func NewAccountRepository() *AccountRepository {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &AccountRepository{
		client: dynamodb.NewFromConfig(cfg),
	}
}

func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	// Ensure account ID is set (use UUID if not set already)
	if account.ID == "" {
		account.ID = uuid.New().String()  // Generate UUID if not provided
	}

	if account.AccountType == "" {
		account.AccountType = "checking"
	}
	// Set the creation timestamp
	account.CreatedAt = time.Now()

	// Log the generated UUID for debugging purposes
	fmt.Printf("Creating account with ID: %s\n", account.ID)

	// Marshal the account struct to a map for DynamoDB
	item, err := attributevalue.MarshalMap(account)
	if err != nil {
		return fmt.Errorf("failed to marshal account: %w", err)
	}

	// Insert the account into DynamoDB
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(AccountsTable),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to insert item into DynamoDB: %w", err)
	}

	// Log the successful insertion (optional)
	fmt.Println("Account created successfully")
	return nil
}


func (r *AccountRepository) GetByID(ctx context.Context, id string) (*models.Account, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(AccountsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var account models.Account
	err = attributevalue.UnmarshalMap(result.Item, &account)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal account: %w", err)
	}

	return &account, nil
}

func (r *AccountRepository) Update(ctx context.Context, account *models.Account) error {
    // Set the updated timestamp
    account.UpdatedAt = time.Now()

    // Update the account's attributes in DynamoDB
    _, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
        TableName: aws.String(AccountsTable),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: account.ID},
        },
        UpdateExpression: aws.String("SET balance = :balance, account_type = :account_type, updated_at = :updated_at"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":balance":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", account.Balance)},
            ":account_type": &types.AttributeValueMemberS{Value: account.AccountType},
            ":updated_at":   &types.AttributeValueMemberS{Value: account.UpdatedAt.Format(time.RFC3339)},
        },
        ReturnValues: types.ReturnValueAllNew,
    })
    if err != nil {
        return fmt.Errorf("failed to update account: %w", err)
    }

    return nil
}

// func (r *AccountRepository) Update(ctx context.Context, account *models.Account) error {
//     // Set the updated timestamp
//     account.UpdatedAt = time.Now()

//     // Marshal the account struct to a map for DynamoDB
//     item, err := attributevalue.MarshalMap(account)
//     if err != nil {
//         return fmt.Errorf("failed to marshal account: %w", err)
//     }

//     // Update the account in DynamoDB
//     _, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
//         TableName: aws.String(AccountsTable),
//         Item:      item,
//     })
//     if err != nil {
//         return fmt.Errorf("failed to update item in DynamoDB: %w", err)
//     }

//     return nil
// }


func (r *AccountRepository) Delete(ctx context.Context, id string) error {
    // Delete the account from DynamoDB
    _, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
        TableName: aws.String(AccountsTable),
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: id},
        },
    })
    if err != nil {
        return fmt.Errorf("failed to delete item from DynamoDB: %w", err)
    }

    return nil
}

// func (r *AccountRepository) UpdateBalance(ctx context.Context, id string, amount float64) error {
// 	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
// 		TableName: aws.String(AccountsTable),
// 		Key: map[string]types.AttributeValue{
// 			"id": &types.AttributeValueMemberS{Value: id},
// 		},
// 		UpdateExpression: aws.String("SET balance = balance + :amount"),
// 		ExpressionAttributeValues: map[string]types.AttributeValue{
// 			":amount": &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", amount)},
// 		},
// 		ReturnValues: types.ReturnValueUpdatedNew,
// 	})
// 	return err
// }

func (r *AccountRepository) ListAll(ctx context.Context) ([]models.Account, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(AccountsTable),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan accounts: %w", err)
	}

	var accounts []models.Account
	err = attributevalue.UnmarshalListOfMaps(result.Items, &accounts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal accounts: %w", err)
	}

	return accounts, nil
}