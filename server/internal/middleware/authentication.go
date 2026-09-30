package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aprimr/tickr/internal/domain"
	"github.com/aprimr/tickr/internal/pkg/response"
	"github.com/aprimr/tickr/internal/utils/jwt"
	"github.com/google/uuid"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyRole   contextKey = "role"
)

func Authenticate(allowedRoles ...domain.UserRole) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "missing authorization header", nil)
				return
			}

			parts := strings.Split(authHeader, " ")
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

			// Check role (if any roles are specified)
			// If no roles, any authenticated user is allowed
			if len(allowedRoles) > 0 {
				roleAllowed := false
				for _, role := range allowedRoles {
					if claims.UserRole == string(role) {
						roleAllowed = true
						break
					}
				}

				if !roleAllowed {
					response.Error(w, http.StatusForbidden, "insufficient permissions", nil)
					return
				}
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
