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

type EquipmentHandler struct {
	svc *service.EquipmentService
}

func NewEquipmentHandler(svc *service.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{svc: svc}
}

func (h *EquipmentHandler) List(w http.ResponseWriter, r *http.Request) {
	page, perPage := ParsePagination(r)
	q := r.URL.Query()

	filter := model.EquipmentFilter{
		Search:  q.Get("search"),
		Page:    page,
		PerPage: perPage,
	}

	if v := q.Get("category_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.CategoryID = &id
		}
	}
	if v := q.Get("department_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.DepartmentID = &id
		}
	}
	if v := q.Get("status"); v != "" {
		s := model.EquipmentStatus(v)
		filter.Status = &s
	}

	list, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WritePaginated(w, list, page, perPage, total)
}

func (h *EquipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	var req model.CreateEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	e, err := h.svc.Create(r.Context(), actorID, &req)
	if err != nil {
		WriteError(w, http.StatusConflict, "CONFLICT", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, e)
}

func (h *EquipmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
		return
	}
	e, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "equipment not found")
		return
	}
	WriteJSON(w, http.StatusOK, e)
}

func (h *EquipmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
		return
	}

	var req model.UpdateEquipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	e, err := h.svc.Update(r.Context(), actorID, id, &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, e)
}

func (h *EquipmentHandler) Archive(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
		return
	}
	if err := h.svc.Archive(r.Context(), actorID, id); err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EquipmentHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to parse multipart form")
		return
	}

	file, fh, err := r.FormFile("photo")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "photo field is required")
		return
	}
	defer file.Close()

	photo, err := h.svc.UploadPhoto(r.Context(), actorID, id, fh)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, photo)
}

func (h *EquipmentHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
		return
	}
	photoID, err := uuid.Parse(chi.URLParam(r, "photoId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid photo id")
		return
	}

	if err := h.svc.DeletePhoto(r.Context(), actorID, id, photoID); err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EquipmentHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="equipment.csv"`)
	w.Write([]byte("\xef\xbb\xbf")) // BOM для корректного открытия в Excel
	if err := h.svc.ExportCSV(r.Context(), w); err != nil {
		return
	}
}
