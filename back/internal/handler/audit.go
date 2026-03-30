package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/service"
)

type AuditHandler struct {
	svc *service.AuditService
}

func NewAuditHandler(svc *service.AuditService) *AuditHandler {
	return &AuditHandler{svc: svc}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := ParsePagination(r)
	q := r.URL.Query()

	filter := model.AuditFilter{
		Action:     q.Get("action"),
		EntityType: q.Get("entity_type"),
		Page:       page,
		PerPage:    perPage,
	}

	if v := q.Get("user_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.UserID = &id
		}
	}
	if v := q.Get("date_from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			filter.DateFrom = &t
		}
	}
	if v := q.Get("date_to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			filter.DateTo = &t
		}
	}

	logs, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WritePaginated(w, logs, page, perPage, total)
}
