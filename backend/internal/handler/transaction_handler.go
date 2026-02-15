package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/asilingas/fambudg/backend/internal/middleware"
	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
	validator          *validator.Validate
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		validator:          validator.New(),
	}
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse query parameters for filters
	filters := &model.TransactionFilters{
		AccountID:  r.URL.Query().Get("accountId"),
		CategoryID: r.URL.Query().Get("categoryId"),
		Type:       r.URL.Query().Get("type"),
		StartDate:  r.URL.Query().Get("startDate"),
		EndDate:    r.URL.Query().Get("endDate"),
	}

	if isSharedStr := r.URL.Query().Get("isShared"); isSharedStr != "" {
		isShared := isSharedStr == "true"
		filters.IsShared = &isShared
	}

	transactions, err := h.transactionService.GetByUserID(r.Context(), userID, filters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction, err := h.transactionService.Create(r.Context(), userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, transaction)
}

func (h *TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "id")
	if transactionID == "" {
		respondWithError(w, http.StatusBadRequest, "missing transaction ID")
		return
	}

	transaction, err := h.transactionService.GetByID(r.Context(), transactionID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, transaction)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "id")
	if transactionID == "" {
		respondWithError(w, http.StatusBadRequest, "missing transaction ID")
		return
	}

	var req model.UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction, err := h.transactionService.Update(r.Context(), transactionID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, transaction)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "id")
	if transactionID == "" {
		respondWithError(w, http.StatusBadRequest, "missing transaction ID")
		return
	}

	if err := h.transactionService.Delete(r.Context(), transactionID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TransactionHandler) GenerateRecurring(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Optional upTo date from query param, defaults to today
	upTo := time.Now()
	if dateStr := r.URL.Query().Get("upTo"); dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid upTo date format, use YYYY-MM-DD")
			return
		}
		upTo = parsed
	}

	result, err := h.transactionService.GenerateRecurring(r.Context(), userID, upTo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, result)
}
