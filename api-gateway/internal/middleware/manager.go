package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const userIDKey contextKey = "userID"

type Manager struct {
	jwtVerifier *JWTVerifier
}

func NewManager(jwtSecret string) *Manager {

	return &Manager{
		jwtVerifier: NewJWTVerifier(jwtSecret),
	}
}

func (mw *Manager) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenStr := extractToken(r)
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		claims, err := mw.jwtVerifier.Verify(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) {
		return ""
	}

	return authHeader[len(prefix):]
}

func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
