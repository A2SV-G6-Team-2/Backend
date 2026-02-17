package http

import (
	"context"
	"net/http"
	"strings"

	"expense_tracker/infrastructure/auth"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

// JWTAuthMiddleware validates Bearer token for /expenses and /categories; sets user ID in context.
// /api-docs and / are left public (no auth required).
func JWTAuthMiddleware(jwtSvc *auth.JWTService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/expenses") || strings.HasPrefix(path, "/categories") {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == "" {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			userID, err := jwtSvc.Validate(tokenStr)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDContextKey, userID.String())
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// UserIDFromRequest returns the user ID from request context (set by JWTAuthMiddleware).
func UserIDFromRequest(r *http.Request) string {
	v := r.Context().Value(UserIDContextKey)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
