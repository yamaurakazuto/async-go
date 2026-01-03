package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"async-go/internal/model"
	"async-go/internal/repository"
	"async-go/internal/service"
)

// BankHandler provides HTTP handlers for the banking API.
type BankHandler struct {
	service *service.BankService
}

func NewBankHandler(service *service.BankService) *BankHandler {
	return &BankHandler{service: service}
}

func (h *BankHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/login", h.handleLogin)
	mux.HandleFunc("/api/balance", h.handleBalance)
	mux.HandleFunc("/api/transfer", h.handleTransfer)
}

func (h *BankHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	userID, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	respondJSON(w, model.LoginResponse{UserID: userID}, http.StatusOK)
}

func (h *BankHandler) handleBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDRaw := r.URL.Query().Get("user_id")
	if userIDRaw == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDRaw, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	balance, err := h.service.Balance(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch balance", http.StatusInternalServerError)
		return
	}

	respondJSON(w, model.BalanceResponse{UserID: userID, Balance: balance}, http.StatusOK)
}

func (h *BankHandler) handleTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.FromUserID == 0 || req.ToAccountNumber == "" || req.Amount <= 0 {
		http.Error(w, "invalid transfer request", http.StatusBadRequest)
		return
	}

	transferID, err := h.service.Transfer(r.Context(), req.FromUserID, req.ToAccountNumber, req.Amount)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientBalance) {
			http.Error(w, "insufficient balance", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to create transfer", http.StatusInternalServerError)
		return
	}

	respondJSON(w, model.TransferResponse{TransferID: transferID}, http.StatusOK)
}

func respondJSON(w http.ResponseWriter, payload any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
