// package handlers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/corebank-api/internal/models"
// 	"github.com/corebank-api/internal/repository"
// 	"github.com/google/uuid"
// )

// type AccountHandler struct {
// 	repo             *repository.AccountRepository
// 	pythonServiceURL string
// 	httpClient       *http.Client
// }

// func NewAccountHandler(repo *repository.AccountRepository, pythonServiceURL string) *AccountHandler {
// 	return &AccountHandler{
// 		repo:             repo,
// 		pythonServiceURL: pythonServiceURL,
// 		httpClient:       &http.Client{Timeout: 5 * time.Second},
// 	}
// }

// func (h *AccountHandler) HandleAccounts(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	switch r.Method {
// 	case http.MethodGet:
// 		h.listAccounts(w, r)
// 	case http.MethodPost:
// 		h.createAccount(w, r)
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// func (h *AccountHandler) HandleAccountByID(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	id := r.URL.Path[len("/accounts/"):]

// 	switch r.Method {
// 	case http.MethodGet:
// 		h.getAccount(w, r, id)
// 	case http.MethodPut:
// 		h.updateAccount(w, r, id)
// 	case http.MethodDelete:
// 		h.deleteAccount(w, r, id)
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// func (h *AccountHandler) listAccounts(w http.ResponseWriter, r *http.Request) {
// 	accounts, err := h.repo.ListAll(r.Context())
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(accounts)
// }

// func (h *AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
// 	var account models.Account
// 	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Set default values
// 	account.ID = uuid.New().String()
// 	account.Balance = 0.0
// 	account.CreatedAt = time.Now()
// 	if account.AccountType == "" {
// 		account.AccountType = "checking"
// 	}

// 	if err := h.repo.Create(r.Context(), &account); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Create initial transaction
// 	if err := h.createInitialTransaction(account.ID); err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to create initial transaction: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(account)
// }

// func (h *AccountHandler) getAccount(w http.ResponseWriter, r *http.Request, id string) {
// 	account, err := h.repo.GetByID(r.Context(), id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if account == nil {
// 		http.Error(w, "Account not found", http.StatusNotFound)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(account)
// }

// func (h *AccountHandler) updateAccount(w http.ResponseWriter, r *http.Request, id string) {
// 	var account models.Account
// 	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	account.ID = id
// 	if err := h.repo.Update(r.Context(), &account); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(account)
// }

// func (h *AccountHandler) deleteAccount(w http.ResponseWriter, r *http.Request, id string) {
// 	if err := h.repo.Delete(r.Context(), id); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusNoContent)
// }

// func (h *AccountHandler) createInitialTransaction(accountID string) error {
// 	transaction := map[string]interface{}{
// 		"account_id": accountID,
// 		"amount":     1000.0,
// 		"type":       "deposit",
// 		"status":     "completed",
// 	}

// 	jsonData, err := json.Marshal(transaction)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal transaction: %w", err)
// 	}

// 	resp, err := h.httpClient.Post(
// 		h.pythonServiceURL+"/transactions",
// 		"application/json",
// 		bytes.NewBuffer(jsonData),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to call transaction service: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
// 		return fmt.Errorf("transaction service returned status: %d", resp.StatusCode)
// 	}

// 	return nil
// }

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	"time"

	"github.com/corebank-api/internal/models"
	"github.com/corebank-api/internal/repository"
	"github.com/google/uuid"
)

type AccountHandler struct {
	repo             *repository.AccountRepository
	pythonServiceURL string
}

func NewAccountHandler(repo *repository.AccountRepository, pythonServiceURL string) *AccountHandler {
	return &AccountHandler{
		repo:             repo,
		pythonServiceURL: pythonServiceURL,
	}
}

func (h *AccountHandler) HandleAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.listAccounts(w, r)
	case http.MethodPost:
		h.createAccount(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AccountHandler) HandleAccountByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Path[len("/accounts/"):]

	switch r.Method {
	case http.MethodGet:
		h.getAccount(w, r, id)
	case http.MethodPut:
		h.updateAccount(w, r, id)
	case http.MethodDelete:
		h.deleteAccount(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AccountHandler) listAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.repo.ListAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(accounts)
}

func (h *AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log before setting ID
	fmt.Println("Creating account with Owner: ", account.Owner)

	// Auto-generate account ID using UUID
	account.ID = uuid.New().String()

	// Set default balance
	account.Balance = 0.0

	// Set the creation timestamp
	account.CreatedAt = time.Now()

	// Call repository to create the account
	if err := h.repo.Create(r.Context(), &account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call Python transaction service to create an initial deposit transaction
	if err := h.callTransactionService(account.ID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to call transaction service: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the created account
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// func (h *AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
// 	var account models.Account
// 	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// If account ID is empty, generate one
// 	if account.ID == "" {
// 		account.ID = uuid.New().String()
// 	}

// 	log.Println("Creating account with ID:", account.ID)

// 	account.CreatedAt = time.Now()

// 	if err := h.repo.Create(r.Context(), &account); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if err := h.callTransactionService(account.ID); err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to call transaction service: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(account)
// }

func (h *AccountHandler) getAccount(w http.ResponseWriter, r *http.Request, id string) {
	account, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if account == nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) updateAccount(w http.ResponseWriter, r *http.Request, id string) {
	var updatedAccount models.Account
	if err := json.NewDecoder(r.Body).Decode(&updatedAccount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingAccount, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existingAccount == nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	updatedAccount.ID = id
	updatedAccount.CreatedAt = existingAccount.CreatedAt

	// Using Update instead of Create for clarity
	if err := h.repo.Update(r.Context(), &updatedAccount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedAccount)
}

func (h *AccountHandler) deleteAccount(w http.ResponseWriter, r *http.Request, id string) {
	err := h.repo.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete account: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AccountHandler) callTransactionService(accountID string) error {
	transaction := models.Transaction{
		AccountID: accountID,
		Amount:    1000,
		Type:      "deposit",
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	txnJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	resp, err := http.Post(h.pythonServiceURL+"/transactions", "application/json", bytes.NewBuffer(txnJSON))
	if err != nil {
		return fmt.Errorf("failed to call transaction service: %w", err)
	}
	defer resp.Body.Close()

	// Accept any 2xx success code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("transaction service returned status: %d", resp.StatusCode)
	}

	return nil
}