// internal/delivery/http/beneficiary_handler.go
package http

import (
	"GonPay_Backend/internal/domain"
	"GonPay_Backend/internal/usecase"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BeneficiaryHandler struct {
	beneficiaryUseCase *usecase.BeneficiaryUseCase
}

type CreateBeneficiaryRequest struct {
	Name              string             `json:"name" validate:"required"`
	AccountIdentifier string             `json:"account_identifier" validate:"required"`
	AccountType       domain.AccountType `json:"account_type" validate:"required"`
	BankName          string             `json:"bank_name"`
}

type UpdateBeneficiaryRequest struct {
	Name              string             `json:"name" validate:"required"`
	AccountIdentifier string             `json:"account_identifier" validate:"required"`
	AccountType       domain.AccountType `json:"account_type" validate:"required"`
	BankName          string             `json:"bank_name"`
}

func NewBeneficiaryHandler(beneficiaryUseCase *usecase.BeneficiaryUseCase) *BeneficiaryHandler {
	return &BeneficiaryHandler{
		beneficiaryUseCase: beneficiaryUseCase,
	}
}

func (h *BeneficiaryHandler) CreateBeneficiary(w http.ResponseWriter, r *http.Request) {
	var req CreateBeneficiaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	beneficiary, err := h.beneficiaryUseCase.CreateBeneficiary(
		userID,
		req.Name,
		req.AccountIdentifier,
		req.AccountType,
		req.BankName,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, beneficiary)
}

func (h *BeneficiaryHandler) UpdateBeneficiary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid beneficiary ID")
		return
	}

	var req UpdateBeneficiaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	beneficiary, err := h.beneficiaryUseCase.UpdateBeneficiary(
		id,
		userID,
		req.Name,
		req.AccountIdentifier,
		req.AccountType,
		req.BankName,
	)
	if err != nil {
		switch err {
		case domain.ErrInvalidOperation:
			respondWithError(w, http.StatusForbidden, "Cannot modify this beneficiary")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, beneficiary)
}

func (h *BeneficiaryHandler) DeleteBeneficiary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid beneficiary ID")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	if err := h.beneficiaryUseCase.DeleteBeneficiary(id, userID); err != nil {
		switch err {
		case domain.ErrInvalidOperation:
			respondWithError(w, http.StatusForbidden, "Cannot delete this beneficiary")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Beneficiary deleted successfully"})
}
 
func (h *BeneficiaryHandler) GetBeneficiary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid beneficiary ID")
		return
	}

	userID := r.Context().Value("user_id").(int64)

	beneficiary, err := h.beneficiaryUseCase.GetBeneficiary(id, userID)
	if err != nil {
		switch err {
		case domain.ErrInvalidOperation:
			respondWithError(w, http.StatusNotFound, "Beneficiary not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, beneficiary)
}

func (h *BeneficiaryHandler) GetUserBeneficiaries(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	beneficiaries, err := h.beneficiaryUseCase.GetUserBeneficiaries(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, beneficiaries)
}
