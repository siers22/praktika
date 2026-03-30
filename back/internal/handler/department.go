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

type DepartmentHandler struct {
	svc *service.DepartmentService
}

func NewDepartmentHandler(svc *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{svc: svc}
}

func (h *DepartmentHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	var req model.CreateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	d, err := h.svc.Create(r.Context(), actorID, &req)
	if err != nil {
		WriteError(w, http.StatusConflict, "CONFLICT", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, d)
}

func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
		return
	}
	var req model.UpdateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	d, err := h.svc.Update(r.Context(), actorID, id, &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, d)
}

func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
		return
	}
	if err := h.svc.Delete(r.Context(), actorID, id); err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
