package http

import (
	"net/http"
)

// Handler provides HTTP handlers for wallet services
type Handler struct {
	// Add use case dependencies here when implementing HTTP endpoints
	// For example:
	// createTxUseCase usecase.CreateTransactionUseCase
	// sendTxUseCase   usecase.SendTransactionUseCase
}

// NewHandler creates a new HTTP handler
func NewHandler() *Handler {
	return &Handler{}
}

// RegisterRoutes registers HTTP routes
func RegisterRoutes(mux *http.ServeMux) {
	// Add HTTP routes here when implementing HTTP endpoints
	// For example:
	// mux.HandleFunc("/api/v1/transactions", handler.handleCreateTransaction)
	// mux.HandleFunc("/api/v1/transactions/send", handler.handleSendTransaction)
}

// Example handler skeleton (to be implemented):
// func (h *Handler) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
// 	// 1. Parse and validate request
// 	// 2. Call use case
// 	// 3. Format and return response
// }
