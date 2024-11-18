// internal/delivery/http/handler.go
package http

import (
	"encoding/json"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

type Handler struct {
	UserHandler        *UserHandler
	WalletHandler      *WalletHandler
	TransactionHandler *TransactionHandler
}

func NewHandler(userHandler *UserHandler, walletHandler *WalletHandler, transactionHandler *TransactionHandler) *Handler {
	return &Handler{
		UserHandler:        userHandler,
		WalletHandler:      walletHandler,
		TransactionHandler: transactionHandler,
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ResponseError{Message: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
