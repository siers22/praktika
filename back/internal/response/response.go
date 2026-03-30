package response

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/siers22/praktika/back/internal/model"
)

type Envelope struct {
	Data any               `json:"data,omitempty"`
	Meta *model.Pagination `json:"meta,omitempty"`
}

type ErrDetail struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrEnvelope struct {
	Error ErrDetail `json:"error"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Data: data})
}

func Paginated(w http.ResponseWriter, data any, page, perPage, total int) {
	meta := model.NewPagination(page, perPage, total)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Envelope{Data: data, Meta: &meta})
}

func Error(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrEnvelope{Error: ErrDetail{Code: code, Message: message}})
}

func ParsePagination(r *http.Request) (page, perPage int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ = strconv.Atoi(r.URL.Query().Get("per_page"))
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 || perPage > 50 {
		perPage = 20
	}
	return
}
