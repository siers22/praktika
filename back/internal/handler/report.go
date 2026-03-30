package handler

import (
	"net/http"

	"github.com/siers22/praktika/back/internal/service"
)

type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) Summary(w http.ResponseWriter, r *http.Request) {
	report, err := h.svc.Summary(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, report)
}

func (h *ReportHandler) ByDepartment(w http.ResponseWriter, r *http.Request) {
	report, err := h.svc.ByDepartment(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, report)
}

func (h *ReportHandler) Movements(w http.ResponseWriter, r *http.Request) {
	mvHandler := &MovementHandler{svc: nil}
	_ = mvHandler
	// Делегируем в MovementHandler.ListAll
	WriteError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "use /api/v1/movements endpoint")
}

func (h *ReportHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.Dashboard(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, data)
}
