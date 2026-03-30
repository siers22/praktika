package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/middleware"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/service"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := ParsePagination(r)
	users, total, err := h.userSvc.List(r.Context(), page, perPage)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WritePaginated(w, users, page, perPage, total)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := h.userSvc.Create(r.Context(), actorID, &req)
	if err != nil {
		WriteError(w, http.StatusConflict, "CONFLICT", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid user id")
		return
	}
	user, err := h.userSvc.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
		return
	}
	WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid user id")
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := h.userSvc.Update(r.Context(), actorID, id, &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid user id")
		return
	}

	var req model.UpdateUserStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}

	if err := h.userSvc.UpdateStatus(r.Context(), actorID, id, req.IsActive); err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
