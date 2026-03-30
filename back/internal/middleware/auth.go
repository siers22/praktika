package middleware

import (
	"net/http"
	"strings"

	"github.com/siers22/praktika/back/internal/model"
	"github.com/siers22/praktika/back/internal/response"
	"github.com/siers22/praktika/back/internal/service"
)

type AuthMiddleware struct {
	authSvc *service.AuthService
}

func NewAuthMiddleware(authSvc *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authSvc: authSvc}
}

// Authenticate validates the JWT and sets user context.
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header format")
			return
		}

		claims, err := m.authSvc.ValidateAccessToken(parts[1])
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
			return
		}

		ctx := setUserID(r.Context(), claims.UserID)
		ctx = setUserRole(ctx, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole returns a middleware that checks the user has one of the given roles.
func RequireRole(roles ...model.Role) func(http.Handler) http.Handler {
	allowed := make(map[model.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := GetUserRole(r)
			if !ok {
				response.Error(w, http.StatusForbidden, "FORBIDDEN", "access denied")
				return
			}
			if _, ok := allowed[role]; !ok {
				response.Error(w, http.StatusForbidden, "FORBIDDEN", "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
