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

type AllowanceHandler struct {
	allowanceService *service.AllowanceService
	validator        *validator.Validate
}

func NewAllowanceHandler(allowanceService *service.AllowanceService) *AllowanceHandler {
	return &AllowanceHandler{
		allowanceService: allowanceService,
		validator:        validator.New(),
	}
}

// List returns all allowances for admin, or own allowance for child
func (h *AllowanceHandler) List(w http.ResponseWriter, r *http.Request) {
	userRole := middleware.GetUserRole(r.Context())
	userID := middleware.GetUserID(r.Context())

	if userRole == "child" {
		// Children see only their own allowance
		allowance, err := h.allowanceService.GetByUserID(r.Context(), userID)
		if err != nil {
			respondWithJSON(w, http.StatusOK, []*model.Allowance{})
			return
		}
		respondWithJSON(w, http.StatusOK, []*model.Allowance{allowance})
		return
	}

	allowances, err := h.allowanceService.GetAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, allowances)
}

// Create sets an allowance for a child (admin only)
func (h *AllowanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateAllowanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	allowance, err := h.allowanceService.Create(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, allowance)
}

// Update updates an allowance (admin only)
func (h *AllowanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	allowanceID := chi.URLParam(r, "id")
	if allowanceID == "" {
		respondWithError(w, http.StatusBadRequest, "missing allowance ID")
		return
	}

	var req model.UpdateAllowanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	allowance, err := h.allowanceService.Update(r.Context(), allowanceID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, allowance)
}
