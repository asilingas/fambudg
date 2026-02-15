package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/asilingas/fambudg/backend/internal/middleware"
	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/service"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Default to current month/year
	now := time.Now()
	month := now.Month()
	year := now.Year()

	if m := r.URL.Query().Get("month"); m != "" {
		parsed, err := strconv.Atoi(m)
		if err != nil || parsed < 1 || parsed > 12 {
			respondWithError(w, http.StatusBadRequest, "invalid month")
			return
		}
		month = time.Month(parsed)
	}

	if y := r.URL.Query().Get("year"); y != "" {
		parsed, err := strconv.Atoi(y)
		if err != nil || parsed < 2000 {
			respondWithError(w, http.StatusBadRequest, "invalid year")
			return
		}
		year = parsed
	}

	dashboard, err := h.reportService.GetDashboard(r.Context(), userID, int(month), year)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, dashboard)
}

func (h *ReportHandler) Monthly(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	month, year, err := parseMonthYear(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	summary, err := h.reportService.GetMonthlySummary(r.Context(), userID, month, year)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, summary)
}

func (h *ReportHandler) ByCategory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	month, year, err := parseMonthYear(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	spending, err := h.reportService.GetSpendingByCategory(r.Context(), userID, month, year)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, spending)
}

func (h *ReportHandler) ByMember(w http.ResponseWriter, r *http.Request) {
	month, year, err := parseMonthYear(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	spending, err := h.reportService.GetSpendingByMember(r.Context(), month, year)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, spending)
}

func (h *ReportHandler) Search(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	filters := &model.SearchFilters{
		UserID:      userID,
		Description: r.URL.Query().Get("description"),
		StartDate:   r.URL.Query().Get("startDate"),
		EndDate:     r.URL.Query().Get("endDate"),
		CategoryID:  r.URL.Query().Get("categoryId"),
		AccountID:   r.URL.Query().Get("accountId"),
	}

	if minStr := r.URL.Query().Get("minAmount"); minStr != "" {
		min, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid minAmount")
			return
		}
		filters.MinAmount = &min
	}

	if maxStr := r.URL.Query().Get("maxAmount"); maxStr != "" {
		max, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid maxAmount")
			return
		}
		filters.MaxAmount = &max
	}

	if tagsStr := r.URL.Query().Get("tags"); tagsStr != "" {
		filters.Tags = strings.Split(tagsStr, ",")
	}

	results, err := h.reportService.SearchTransactions(r.Context(), filters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}

func parseMonthYear(r *http.Request) (int, int, error) {
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	if monthStr == "" || yearStr == "" {
		return 0, 0, fmt.Errorf("month and year are required")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("invalid month")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 {
		return 0, 0, fmt.Errorf("invalid year")
	}

	return month, year, nil
}
