package middleware

import (
	"context"
	"fmt"
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

		fmt.Println("token:" + tokenStr)
		claims, err := mw.jwtVerifier.Verify(tokenStr)
		if err != nil {
			fmt.Println("JWT ERROR:", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		fmt.Println("USER ID:", claims.UID)

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
