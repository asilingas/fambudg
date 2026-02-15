package handler

import (
	"encoding/json"
	"net/http"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type SavingGoalHandler struct {
	savingGoalService *service.SavingGoalService
	validator         *validator.Validate
}

func NewSavingGoalHandler(savingGoalService *service.SavingGoalService) *SavingGoalHandler {
	return &SavingGoalHandler{
		savingGoalService: savingGoalService,
		validator:         validator.New(),
	}
}

func (h *SavingGoalHandler) List(w http.ResponseWriter, r *http.Request) {
	goals, err := h.savingGoalService.GetAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, goals)
}

func (h *SavingGoalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateSavingGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	goal, err := h.savingGoalService.Create(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, goal)
}

func (h *SavingGoalHandler) Update(w http.ResponseWriter, r *http.Request) {
	goalID := chi.URLParam(r, "id")
	if goalID == "" {
		respondWithError(w, http.StatusBadRequest, "missing saving goal ID")
		return
	}

	var req model.UpdateSavingGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	goal, err := h.savingGoalService.Update(r.Context(), goalID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, goal)
}

func (h *SavingGoalHandler) Contribute(w http.ResponseWriter, r *http.Request) {
	goalID := chi.URLParam(r, "id")
	if goalID == "" {
		respondWithError(w, http.StatusBadRequest, "missing saving goal ID")
		return
	}

	var req model.ContributeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	goal, err := h.savingGoalService.Contribute(r.Context(), goalID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, goal)
}
