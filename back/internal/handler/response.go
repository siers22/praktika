// Package handler re-exports response helpers for convenience.
// Actual implementation lives in internal/response to avoid import cycles.
package handler

import (
	"net/http"

	"github.com/siers22/praktika/back/internal/response"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	response.JSON(w, status, data)
}

func WritePaginated(w http.ResponseWriter, data any, page, perPage, total int) {
	response.Paginated(w, data, page, perPage, total)
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	response.Error(w, status, code, message)
}

func ParsePagination(r *http.Request) (page, perPage int) {
	return response.ParsePagination(r)
}
