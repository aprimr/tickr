package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aprimr/tickr/internal/pkg/response"
	"github.com/aprimr/tickr/internal/utils/jwt"
	"github.com/google/uuid"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyRole   contextKey = "role"
)

func Authenticate(allowedRole string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "missing authorization header", nil)
				return
			}

			// Check Authorization header format
			parts := strings.Split(authHeader, "")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(w, http.StatusUnauthorized, "invalid authorization header format", nil)
				return
			}

			tokenString := parts[1]

			// Verify access token
			claims, err := jwt.VerifyAccessToken(tokenString)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid or expired access token", nil)
				return
			}

			// Check role
			if claims.UserRole != allowedRole {
				response.Error(w, http.StatusForbidden, "insufficient permissions", nil)
				return
			}

			// Add values to context
			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyRole, claims.UserRole)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper functions
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(uuid.UUID)
	return id, ok
}

func GetRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextKeyRole).(string)
	return role, ok
}
