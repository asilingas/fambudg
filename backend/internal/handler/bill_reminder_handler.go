package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asilingas/fambudg/backend/internal/middleware"
	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type BillReminderHandler struct {
	billReminderService *service.BillReminderService
	validator           *validator.Validate
}

func NewBillReminderHandler(billReminderService *service.BillReminderService) *BillReminderHandler {
	return &BillReminderHandler{
		billReminderService: billReminderService,
		validator:           validator.New(),
	}
}

func (h *BillReminderHandler) List(w http.ResponseWriter, r *http.Request) {
	bills, err := h.billReminderService.GetAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, bills)
}

func (h *BillReminderHandler) Upcoming(w http.ResponseWriter, r *http.Request) {
	days := 30
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		parsed, err := strconv.Atoi(daysStr)
		if err != nil || parsed < 1 {
			respondWithError(w, http.StatusBadRequest, "invalid days parameter")
			return
		}
		days = parsed
	}

	bills, err := h.billReminderService.GetUpcoming(r.Context(), days)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, bills)
}

func (h *BillReminderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBillReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	bill, err := h.billReminderService.Create(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, bill)
}

func (h *BillReminderHandler) Update(w http.ResponseWriter, r *http.Request) {
	billID := chi.URLParam(r, "id")
	if billID == "" {
		respondWithError(w, http.StatusBadRequest, "missing bill reminder ID")
		return
	}

	var req model.UpdateBillReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	bill, err := h.billReminderService.Update(r.Context(), billID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, bill)
}

func (h *BillReminderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	billID := chi.URLParam(r, "id")
	if billID == "" {
		respondWithError(w, http.StatusBadRequest, "missing bill reminder ID")
		return
	}

	if err := h.billReminderService.Delete(r.Context(), billID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BillReminderHandler) Pay(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	billID := chi.URLParam(r, "id")
	if billID == "" {
		respondWithError(w, http.StatusBadRequest, "missing bill reminder ID")
		return
	}

	var req model.PayBillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	transaction, err := h.billReminderService.Pay(r.Context(), userID, billID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, transaction)
}
