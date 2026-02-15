package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/asilingas/fambudg/backend/internal/middleware"
	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/service"
)

type ImportExportHandler struct {
	transactionService *service.TransactionService
}

func NewImportExportHandler(transactionService *service.TransactionService) *ImportExportHandler {
	return &ImportExportHandler{transactionService: transactionService}
}

func (h *ImportExportHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	filters := &model.TransactionFilters{
		StartDate: r.URL.Query().Get("startDate"),
		EndDate:   r.URL.Query().Get("endDate"),
	}

	transactions, err := h.transactionService.GetByUserID(r.Context(), userID, filters)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=transactions.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Header row
	writer.Write([]string{"date", "amount", "type", "description", "category_id", "account_id", "is_shared"})

	for _, t := range transactions {
		writer.Write([]string{
			t.Date.Format("2006-01-02"),
			strconv.FormatInt(t.Amount, 10),
			t.Type,
			t.Description,
			t.CategoryID,
			t.AccountID,
			strconv.FormatBool(t.IsShared),
		})
	}
}

func (h *ImportExportHandler) ImportCSV(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "missing CSV file")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and skip header
	_, err = reader.Read()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to read CSV header")
		return
	}

	var imported int
	var errors []string

	records, err := reader.ReadAll()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse CSV")
		return
	}

	for i, record := range records {
		if len(record) < 6 {
			errors = append(errors, fmt.Sprintf("row %d: insufficient columns", i+2))
			continue
		}

		amount, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			errors = append(errors, fmt.Sprintf("row %d: invalid amount", i+2))
			continue
		}

		isShared := true
		if len(record) > 6 {
			isShared, _ = strconv.ParseBool(record[6])
		}

		req := &model.CreateTransactionRequest{
			Date:        record[0],
			Amount:      amount,
			Type:        record[2],
			Description: record[3],
			CategoryID:  record[4],
			AccountID:   record[5],
			IsShared:    isShared,
		}

		_, err = h.transactionService.Create(r.Context(), userID, req)
		if err != nil {
			errors = append(errors, fmt.Sprintf("row %d: %s", i+2, err.Error()))
			continue
		}

		imported++
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"imported": imported,
		"errors":   errors,
	})
}
