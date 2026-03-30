package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/middleware"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/service"
)

type MovementHandler struct {
	svc *service.MovementService
}

func NewMovementHandler(svc *service.MovementService) *MovementHandler {
	return &MovementHandler{svc: svc}
}

func (h *MovementHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	var req model.CreateMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	mv, err := h.svc.Create(r.Context(), actorID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, mv)
}

func (h *MovementHandler) ListByEquipment(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid equipment id")
		return
	}
	page, perPage := ParsePagination(r)
	list, total, err := h.svc.ListByEquipment(r.Context(), id, page, perPage)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WritePaginated(w, list, page, perPage, total)
}

func (h *MovementHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	page, perPage := ParsePagination(r)
	q := r.URL.Query()

	var dateFrom, dateTo *time.Time
	if v := q.Get("date_from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			dateFrom = &t
		}
	}
	if v := q.Get("date_to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			dateTo = &t
		}
	}

	list, total, err := h.svc.ListAll(r.Context(), dateFrom, dateTo, page, perPage)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WritePaginated(w, list, page, perPage, total)
}
