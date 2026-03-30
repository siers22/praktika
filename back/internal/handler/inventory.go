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

type InventoryHandler struct {
	svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) ListSessions(w http.ResponseWriter, r *http.Request) {
	page, perPage := ParsePagination(r)
	list, total, err := h.svc.ListSessions(r.Context(), page, perPage)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WritePaginated(w, list, page, perPage, total)
}

func (h *InventoryHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	var req model.CreateInventorySessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	session, err := h.svc.CreateSession(r.Context(), actorID, &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, session)
}

func (h *InventoryHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
		return
	}
	session, err := h.svc.GetSession(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "session not found")
		return
	}
	items, err := h.svc.GetItems(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{"session": session, "items": items})
}

func (h *InventoryHandler) CheckItem(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	sessionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid session id")
		return
	}

	var req model.CheckInventoryItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	item, err := h.svc.CheckItem(r.Context(), actorID, sessionID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, item)
}

func (h *InventoryHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	sessionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid session id")
		return
	}
	itemID, err := uuid.Parse(chi.URLParam(r, "itemId"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid item id")
		return
	}

	var req model.UpdateInventoryItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON body")
		return
	}
	if err := validate.Struct(req); err != nil {
		WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	item, err := h.svc.UpdateItem(r.Context(), actorID, sessionID, itemID, &req)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, item)
}

func (h *InventoryHandler) CompleteSession(w http.ResponseWriter, r *http.Request) {
	actorID, _ := middleware.GetUserID(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
		return
	}
	if err := h.svc.CompleteSession(r.Context(), actorID, id); err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *InventoryHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid id")
		return
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="inventory.csv"`)
	w.Write([]byte("\xef\xbb\xbf"))
	_ = h.svc.ExportCSV(r.Context(), id, w)
}
