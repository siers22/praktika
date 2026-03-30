package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/siers22/praktika/back/internal/middleware"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/service"
)

var validate = validator.New()

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	pair, err := h.authSvc.Login(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			WriteError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid username or password")
		case errors.Is(err, service.ErrUserInactive):
			WriteError(w, http.StatusForbidden, "USER_INACTIVE", "account is disabled")
		default:
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	WriteJSON(w, http.StatusOK, pair)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	pair, err := h.authSvc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "INVALID_TOKEN", "invalid or expired refresh token")
		return
	}

	WriteJSON(w, http.StatusOK, pair)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	_ = h.authSvc.Logout(r.Context(), req.RefreshToken)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var req model.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if err := h.authSvc.ChangePassword(r.Context(), userID, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
