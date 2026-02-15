package handler

import (
	"encoding/json"
	"net/http"

	"github.com/asilingas/fambudg/backend/internal/middleware"
	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type AccountHandler struct {
	accountService *service.AccountService
	validator      *validator.Validate
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		validator:      validator.New(),
	}
}

// List returns accounts — admin sees all, others see own
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role := middleware.GetUserRole(r.Context())

	var accounts []*model.Account
	var err error

	if role == "admin" {
		accounts, err = h.accountService.GetAll(r.Context())
	} else {
		accounts, err = h.accountService.GetByUserID(r.Context(), userID)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, accounts)
}

// Create creates a new account
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.accountService.Create(r.Context(), userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, account)
}

// Get returns a single account by ID with ownership check
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role := middleware.GetUserRole(r.Context())

	accountID := chi.URLParam(r, "id")
	if accountID == "" {
		respondWithError(w, http.StatusBadRequest, "missing account ID")
		return
	}

	account, err := h.accountService.GetByID(r.Context(), accountID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// Non-admin users can only see their own accounts
	if role != "admin" && account.UserID != userID {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	respondWithJSON(w, http.StatusOK, account)
}

// Update updates an existing account with ownership check
func (h *AccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role := middleware.GetUserRole(r.Context())

	accountID := chi.URLParam(r, "id")
	if accountID == "" {
		respondWithError(w, http.StatusBadRequest, "missing account ID")
		return
	}

	// Check ownership for non-admin
	existing, err := h.accountService.GetByID(r.Context(), accountID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if role != "admin" && existing.UserID != userID {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req model.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.accountService.Update(r.Context(), accountID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, account)
}

// Delete deletes an account with ownership check
func (h *AccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role := middleware.GetUserRole(r.Context())

	accountID := chi.URLParam(r, "id")
	if accountID == "" {
		respondWithError(w, http.StatusBadRequest, "missing account ID")
		return
	}

	// Check ownership for non-admin
	existing, err := h.accountService.GetByID(r.Context(), accountID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	if role != "admin" && existing.UserID != userID {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	if err := h.accountService.Delete(r.Context(), accountID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
