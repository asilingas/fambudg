package handler

import (
	"encoding/json"
	"net/http"

	"github.com/asilingas/fambudg/backend/internal/middleware"
	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/service"
	"github.com/go-playground/validator/v10"
)

type TransferHandler struct {
	transactionService *service.TransactionService
	validator          *validator.Validate
}

func NewTransferHandler(transactionService *service.TransactionService) *TransferHandler {
	return &TransferHandler{
		transactionService: transactionService,
		validator:          validator.New(),
	}
}

func (h *TransferHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.FromAccountID == req.ToAccountID {
		respondWithError(w, http.StatusBadRequest, "cannot transfer to the same account")
		return
	}

	toAccountID := req.ToAccountID
	txReq := &model.CreateTransactionRequest{
		AccountID:           req.FromAccountID,
		CategoryID:          "", // transfers don't need a category
		Amount:              -req.Amount,
		Type:                "transfer",
		Description:         req.Description,
		Date:                req.Date,
		IsShared:            true,
		TransferToAccountID: &toAccountID,
	}

	transaction, err := h.transactionService.Create(r.Context(), userID, txReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, transaction)
}
