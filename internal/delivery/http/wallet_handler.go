// internal/delivery/http/wallet_handler.go
package http

import (
	"GonPay_Backend/internal/domain"
	"GonPay_Backend/internal/usecase"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type WalletHandler struct {
	walletUseCase *usecase.WalletUseCase
}

type TransferRequest struct {
	SourceWalletID      int64   `json:"source_wallet_id" validate:"required"`
	DestinationWalletID int64   `json:"destination_wallet_id" validate:"required"`
	Amount              float64 `json:"amount" validate:"required,gt=0"`
	Description         string  `json:"description"`
}

type MoneyRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description"`
}

func NewWalletHandler(walletUseCase *usecase.WalletUseCase) *WalletHandler {
	return &WalletHandler{
		walletUseCase: walletUseCase,
	}
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	wallet, err := h.walletUseCase.CreateWallet(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, wallet)
}

func (h *WalletHandler) GetUserWallets(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	wallets, err := h.walletUseCase.GetUserWallets(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, wallets)
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	wallet, err := h.walletUseCase.GetWallet(walletID)
	if err != nil {
		switch err {
		case domain.ErrWalletNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, wallet)
}

func (h *WalletHandler) DeactivateWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	if err := h.walletUseCase.DeactivateWallet(walletID); err != nil {
		switch err {
		case domain.ErrWalletNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Wallet deactivated successfully"})
}

func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tx, err := h.walletUseCase.Transfer(req.SourceWalletID, req.DestinationWalletID, req.Amount)
	if err != nil {
		switch err {
		case domain.ErrInsufficientFunds:
			respondWithError(w, http.StatusBadRequest, err.Error())
		case domain.ErrWalletNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, tx)
}

func (h *WalletHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	var req MoneyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tx, err := h.walletUseCase.Deposit(walletID, req.Amount)
	if err != nil {
		switch err {
		case domain.ErrWalletNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, tx)
}

func (h *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	var req MoneyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	tx, err := h.walletUseCase.Withdraw(walletID, req.Amount)
	if err != nil {
		switch err {
		case domain.ErrInsufficientFunds:
			respondWithError(w, http.StatusBadRequest, err.Error())
		case domain.ErrWalletNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, tx)
}
