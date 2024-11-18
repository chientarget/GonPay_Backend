// internal/delivery/http/payment_method_handler.go
package http

import (
	"GonPay_Backend/internal/domain"
	"GonPay_Backend/internal/usecase"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type PaymentMethodHandler struct {
	paymentMethodUseCase *usecase.PaymentMethodUseCase
}

type CreatePaymentMethodRequest struct {
	MethodType    domain.PaymentMethodType `json:"method_type" validate:"required"`
	AccountNumber string                   `json:"account_number" validate:"required"`
	BankName      string                   `json:"bank_name"`
	IsDefault     bool                     `json:"is_default"`
}

type UpdatePaymentMethodRequest struct {
	MethodType    domain.PaymentMethodType `json:"method_type" validate:"required"`
	AccountNumber string                   `json:"account_number" validate:"required"`
	BankName      string                   `json:"bank_name"`
	IsDefault     bool                     `json:"is_default"`
}

func NewPaymentMethodHandler(paymentMethodUseCase *usecase.PaymentMethodUseCase) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		paymentMethodUseCase: paymentMethodUseCase,
	}
}

func (h *PaymentMethodHandler) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	paymentMethod, err := h.paymentMethodUseCase.CreatePaymentMethod(
		userID,
		req.MethodType,
		req.AccountNumber,
		req.BankName,
		req.IsDefault,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, paymentMethod)
}

// internal/delivery/http/payment_method_handler.go

func (h *PaymentMethodHandler) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payment method ID")
		return
	}

	var req UpdatePaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	paymentMethod, err := h.paymentMethodUseCase.UpdatePaymentMethod(
		id,
		userID,
		req.MethodType,
		req.AccountNumber,
		req.BankName,
		req.IsDefault,
	)
	if err != nil {
		switch err {
		case domain.ErrInvalidOperation:
			respondWithError(w, http.StatusForbidden, "Cannot modify this payment method")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, paymentMethod)
}

func (h *PaymentMethodHandler) DeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payment method ID")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	if err := h.paymentMethodUseCase.DeletePaymentMethod(id, userID); err != nil {
		switch err {
		case domain.ErrInvalidOperation:
			respondWithError(w, http.StatusForbidden, "Cannot delete this payment method")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Payment method deleted successfully"})
}

func (h *PaymentMethodHandler) GetPaymentMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payment method ID")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	paymentMethod, err := h.paymentMethodUseCase.GetPaymentMethod(id, userID)
	if err != nil {
		switch err {
		case domain.ErrInvalidOperation:
			respondWithError(w, http.StatusNotFound, "Payment method not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, paymentMethod)
}

func (h *PaymentMethodHandler) GetUserPaymentMethods(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	paymentMethods, err := h.paymentMethodUseCase.GetUserPaymentMethods(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, paymentMethods)
}

func (h *PaymentMethodHandler) SetDefaultPaymentMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payment method ID")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	if err := h.paymentMethodUseCase.SetDefaultPaymentMethod(id, userID); err != nil {
		switch err {
		case domain.ErrInvalidOperation:
			respondWithError(w, http.StatusForbidden, "Cannot set this payment method as default")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Payment method set as default successfully"})
}
