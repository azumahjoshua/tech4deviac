package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/corebank-api/internal/models"
	"github.com/corebank-api/internal/repository"
)

type TransactionHandler struct {
	accountRepo      *repository.AccountRepository
	pythonServiceURL string
	httpClient       *http.Client
}

func NewTransactionHandler(
	accountRepo *repository.AccountRepository,
	pythonServiceURL string,
) *TransactionHandler {
	return &TransactionHandler{
		accountRepo:      accountRepo,
		pythonServiceURL: pythonServiceURL,
		httpClient:       &http.Client{Timeout: 5 * time.Second},
	}
}

func (h *TransactionHandler) HandleTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// For POST requests, verify account exists first
	if r.Method == http.MethodPost {
		var txn models.Transaction
		if err := json.NewDecoder(r.Body).Decode(&txn); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Verify account exists in Go's database
		account, err := h.accountRepo.GetByID(r.Context(), txn.AccountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if account == nil {
			http.Error(w, "Account not found", http.StatusBadRequest)
			return
		}

		// Marshal the transaction back to JSON for forwarding
		txnBytes, err := json.Marshal(txn)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal transaction: %v", err), http.StatusInternalServerError)
			return
		}

		// Reset the body for forwarding
		r.Body = io.NopCloser(bytes.NewBuffer(txnBytes))
	}

	// Use path.Join to avoid slash issues in URL construction
	targetURL, err := url.JoinPath(h.pythonServiceURL, r.URL.Path)
	if err != nil {
		http.Error(w, "invalid URL path", http.StatusInternalServerError)
		return
	}

	// Forward the request to the Python service
	forwardReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create forward request: %v", err), http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	forwardReq.Header = r.Header

	// Execute the forward request
	resp, err := h.httpClient.Do(forwardReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to forward request to transaction service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response from the Python service back to the client
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (h *TransactionHandler) HandleTransactionByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract transaction ID from the URL path
	txnID := r.URL.Path[len("/transactions/"):]

	// Construct target URL safely
	targetURL, err := url.JoinPath(h.pythonServiceURL, "/transactions", txnID)
	if err != nil {
		http.Error(w, "invalid URL path", http.StatusInternalServerError)
		return
	}

	// Forward the request to the Python service
	forwardReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create forward request: %v", err), http.StatusInternalServerError)
		return
	}

	forwardReq.Header = r.Header

	// Execute the forward request
	resp, err := h.httpClient.Do(forwardReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to forward request to transaction service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response from the Python service back to the client
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (h *TransactionHandler) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Construct target URL safely
	targetURL, err := url.JoinPath(h.pythonServiceURL, "/transactions")
	if err != nil {
		http.Error(w, "invalid URL path", http.StatusInternalServerError)
		return
	}

	// Forward the request to the Python service
	forwardReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create forward request: %v", err), http.StatusInternalServerError)
		return
	}

	forwardReq.Header = r.Header

	// Execute the forward request
	resp, err := h.httpClient.Do(forwardReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to forward request to transaction service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response from the Python service back to the client
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

// package handlers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"time"

// 	"github.com/corebank-api/internal/models"
// 	"github.com/corebank-api/internal/repository"
// )

// type TransactionHandler struct {
// 	accountRepo      *repository.AccountRepository
// 	pythonServiceURL string
// 	httpClient       *http.Client
// }

// func NewTransactionHandler(
// 	accountRepo *repository.AccountRepository,
// 	pythonServiceURL string,
// ) *TransactionHandler {
// 	return &TransactionHandler{
// 		accountRepo:      accountRepo,
// 		pythonServiceURL: pythonServiceURL,
// 		httpClient:       &http.Client{Timeout: 5 * time.Second},
// 	}
// }
// func (h *TransactionHandler) HandleTransactions(w http.ResponseWriter, r *http.Request) {
//     w.Header().Set("Content-Type", "application/json")

//     // For POST requests, verify account exists first
//     if r.Method == http.MethodPost {
//         var txn models.Transaction
//         if err := json.NewDecoder(r.Body).Decode(&txn); err != nil {
//             http.Error(w, err.Error(), http.StatusBadRequest)
//             return
//         }
        
//         // Verify account exists in Go's database
//         account, err := h.accountRepo.GetByID(r.Context(), txn.AccountID)
//         if err != nil {
//             http.Error(w, err.Error(), http.StatusInternalServerError)
//             return
//         }
//         if account == nil {
//             http.Error(w, "Account not found", http.StatusBadRequest)
//             return
//         }

//         // Marshal the transaction back to JSON for forwarding
//         txnBytes, err := json.Marshal(txn)
//         if err != nil {
//             http.Error(w, fmt.Sprintf("failed to marshal transaction: %v", err), http.StatusInternalServerError)
//             return
//         }
        
//         // Reset the body for forwarding
//         r.Body = io.NopCloser(bytes.NewBuffer(txnBytes))
//     }

//     // Use path.Join to avoid slash issues in URL construction
//     targetURL, err := url.JoinPath(h.pythonServiceURL, r.URL.Path)
//     if err != nil {
//         http.Error(w, "invalid URL path", http.StatusInternalServerError)
//         return
//     }

//     // Forward the request to the Python service
//     forwardReq, err := http.NewRequest(r.Method, targetURL, r.Body)
//     if err != nil {
//         http.Error(w, fmt.Sprintf("failed to create forward request: %v", err), http.StatusInternalServerError)
//         return
//     }

//     // Copy headers from the original request
//     forwardReq.Header = r.Header

//     // Execute the forward request
//     resp, err := h.httpClient.Do(forwardReq)
//     if err != nil {
//         http.Error(w, fmt.Sprintf("failed to forward request to transaction service: %v", err), http.StatusInternalServerError)
//         return
//     }
//     defer resp.Body.Close()

//     // Copy the response from the Python service back to the client
//     w.WriteHeader(resp.StatusCode)
//     if _, err := io.Copy(w, resp.Body); err != nil {
//         log.Printf("failed to write response: %v", err)
//     }
// }
// // func (h *TransactionHandler) HandleTransactions(w http.ResponseWriter, r *http.Request) {
// // 	w.Header().Set("Content-Type", "application/json")

// // 	// For POST requests, verify account exists first
// // 	if r.Method == http.MethodPost {
// // 		var txn models.Transaction
// // 		if err := json.NewDecoder(r.Body).Decode(&txn); err != nil {
// // 			http.Error(w, err.Error(), http.StatusBadRequest)
// // 			return
// // 		}

// // 		// Reset the body so it can be read again
// // 		// Use json.Marshal to convert txn to a JSON byte slice
// // 		txnBytes, err := json.Marshal(txn)
// // 		if err != nil {
// // 			http.Error(w, fmt.Sprintf("failed to marshal transaction: %v", err), http.StatusInternalServerError)
// // 			return
// // 		}

// // 		// Reset the body with the marshaled JSON
// // 		r.Body = io.NopCloser(bytes.NewBuffer(txnBytes))

// // 		// Verify account exists in Go's database
// // 		account, err := h.accountRepo.GetByID(r.Context(), txn.AccountID)
// // 		if err != nil {
// // 			http.Error(w, err.Error(), http.StatusInternalServerError)
// // 			return
// // 		}
// // 		if account == nil {
// // 			http.Error(w, "Account not found", http.StatusBadRequest)
// // 			return
// // 		}
// // 	}

// // 	// Forward the request to the Python service
// // 	forwardReq, err := http.NewRequest(r.Method, h.pythonServiceURL+r.URL.Path, r.Body)
// // 	if err != nil {
// // 		http.Error(w, fmt.Sprintf("failed to create forward request: %v", err), http.StatusInternalServerError)
// // 		return
// // 	}

// // 	// Copy headers from the original request
// // 	forwardReq.Header = r.Header

// // 	// Execute the forward request
// // 	resp, err := h.httpClient.Do(forwardReq)
// // 	if err != nil {
// // 		http.Error(w, fmt.Sprintf("failed to forward request to transaction service: %v", err), http.StatusInternalServerError)
// // 		return
// // 	}
// // 	defer resp.Body.Close()

// // 	// Copy the response from the Python service back to the client
// // 	w.WriteHeader(resp.StatusCode)
// // 	io.Copy(w, resp.Body)
// // }

// // HandleTransactionByID - Handle the GET request for a transaction by ID
// func (h *TransactionHandler) HandleTransactionByID(w http.ResponseWriter, r *http.Request) {
//     w.Header().Set("Content-Type", "application/json")

//     // Extract transaction ID from the URL path
//     txnID := r.URL.Path[len("/transactions/"):]

//     // Forward the request to the Python service
//     forwardReq, err := http.NewRequest(r.Method, fmt.Sprintf("%s/transactions/%s", h.pythonServiceURL, txnID), r.Body)
//     if err != nil {
//         http.Error(w, fmt.Sprintf("failed to create forward request: %v", err), http.StatusInternalServerError)
//         return
//     }

//     forwardReq.Header = r.Header

//     // Execute the forward request
//     resp, err := h.httpClient.Do(forwardReq)
//     if err != nil {
//         http.Error(w, fmt.Sprintf("failed to forward request to transaction service: %v", err), http.StatusInternalServerError)
//         return
//     }
//     defer resp.Body.Close()

//     // Copy the response from the Python service back to the client
//     w.WriteHeader(resp.StatusCode)
//     io.Copy(w, resp.Body)
// }

// // HandleGetTransactions - Handle the GET request for listing all transactions
// func (h *TransactionHandler) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
//     w.Header().Set("Content-Type", "application/json")

//     // Forward the request to the Python service
//     forwardReq, err := http.NewRequest(r.Method, fmt.Sprintf("%s/transactions", h.pythonServiceURL), r.Body)
//     if err != nil {
//         http.Error(w, fmt.Sprintf("failed to create forward request: %v", err), http.StatusInternalServerError)
//         return
//     }

//     forwardReq.Header = r.Header

//     // Execute the forward request
//     resp, err := h.httpClient.Do(forwardReq)
//     if err != nil {
//         http.Error(w, fmt.Sprintf("failed to forward request to transaction service: %v", err), http.StatusInternalServerError)
//         return
//     }
//     defer resp.Body.Close()

//     // Copy the response from the Python service back to the client
//     w.WriteHeader(resp.StatusCode)
//     io.Copy(w, resp.Body)
// }