package http

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

// MockUserIDMiddleware reads X-User-ID header and sets user ID in context (for Team 2 mock users)
func MockUserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromRequest returns the user ID from request context (set by MockUserIDMiddleware)
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
