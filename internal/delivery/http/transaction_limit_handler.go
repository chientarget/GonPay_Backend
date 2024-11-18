// internal/delivery/http/transaction_limit_handler.go
package http

import (
	"GonPay_Backend/internal/domain"
	"GonPay_Backend/internal/usecase"
	"encoding/json"
	"net/http"
)

type TransactionLimitHandler struct {
	limitUseCase *usecase.TransactionLimitUseCase
}

type SetLimitRequest struct {
	TransactionType domain.TransactionType `json:"transaction_type" validate:"required"`
	DailyLimit      float64                `json:"daily_limit" validate:"required,gt=0"`
	MonthlyLimit    float64                `json:"monthly_limit" validate:"required,gt=0"`
}

func NewTransactionLimitHandler(limitUseCase *usecase.TransactionLimitUseCase) *TransactionLimitHandler {
	return &TransactionLimitHandler{
		limitUseCase: limitUseCase,
	}
}

func (h *TransactionLimitHandler) SetLimit(w http.ResponseWriter, r *http.Request) {
	var req SetLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	limit, err := h.limitUseCase.SetTransactionLimit(
		userID,
		req.TransactionType,
		req.DailyLimit,
		req.MonthlyLimit,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, limit)
}

func (h *TransactionLimitHandler) GetLimits(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	limits, err := h.limitUseCase.GetUserLimits(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, limits)
}
