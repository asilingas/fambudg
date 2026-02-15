package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type BudgetHandler struct {
	budgetService *service.BudgetService
	validator     *validator.Validate
}

func NewBudgetHandler(budgetService *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
		validator:     validator.New(),
	}
}

func (h *BudgetHandler) List(w http.ResponseWriter, r *http.Request) {
	filters := &model.BudgetFilters{}

	if m := r.URL.Query().Get("month"); m != "" {
		month, err := strconv.Atoi(m)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid month")
			return
		}
		filters.Month = month
	}

	if y := r.URL.Query().Get("year"); y != "" {
		year, err := strconv.Atoi(y)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid year")
			return
		}
		filters.Year = year
	}

	budgets, err := h.budgetService.GetAll(r.Context(), filters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, budgets)
}

func (h *BudgetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	budget, err := h.budgetService.Create(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, budget)
}

func (h *BudgetHandler) Update(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		respondWithError(w, http.StatusBadRequest, "missing budget ID")
		return
	}

	var req model.UpdateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	budget, err := h.budgetService.Update(r.Context(), budgetID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, budget)
}

func (h *BudgetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		respondWithError(w, http.StatusBadRequest, "missing budget ID")
		return
	}

	if err := h.budgetService.Delete(r.Context(), budgetID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BudgetHandler) Summary(w http.ResponseWriter, r *http.Request) {
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	if monthStr == "" || yearStr == "" {
		respondWithError(w, http.StatusBadRequest, "month and year are required")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		respondWithError(w, http.StatusBadRequest, "invalid month")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 {
		respondWithError(w, http.StatusBadRequest, "invalid year")
		return
	}

	summaries, err := h.budgetService.GetSummary(r.Context(), month, year)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, summaries)
}
