// internal/delivery/http/transaction_handler.go
package http

import (
	"GonPay_Backend/internal/usecase"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	transactionUseCase *usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		transactionUseCase: transactionUseCase,
	}
}

func (h *TransactionHandler) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	transactions, err := h.transactionUseCase.GetUserTransactions(userID, page, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Add metadata for pagination
	response := map[string]interface{}{
		"data": transactions,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": len(transactions),
		},
	}

	respondWithJSON(w, http.StatusOK, response)
}
