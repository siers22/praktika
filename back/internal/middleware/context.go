package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/siers22/praktika/back/internal/model"
)

type contextKey string

const (
	ctxUserID   contextKey = "user_id"
	ctxUserRole contextKey = "user_role"
)

func setUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxUserID, id)
}

func setUserRole(ctx context.Context, role model.Role) context.Context {
	return context.WithValue(ctx, ctxUserRole, role)
}

func GetUserID(r *http.Request) (uuid.UUID, bool) {
	id, ok := r.Context().Value(ctxUserID).(uuid.UUID)
	return id, ok
}

func GetUserRole(r *http.Request) (model.Role, bool) {
	role, ok := r.Context().Value(ctxUserRole).(model.Role)
	return role, ok
}
