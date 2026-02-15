package handler

import (
	"encoding/json"
	"net/http"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	authService *service.AuthService
	validator   *validator.Validate
}

func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.authService.ListUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.CreateUser(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	if err := h.authService.DeleteUser(r.Context(), userID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
