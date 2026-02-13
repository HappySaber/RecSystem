package middleware

import (
	"api-gateway/internal/client"
	"context"
	"net/http"
)

type Manager struct {
	ssoClient *client.SSOClient
}

func NewManager(sso *client.SSOClient) *Manager {
	return &Manager{
		ssoClient: sso,
	}
}

func (m *Manager) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := extractToken(r)

		// resp, err := m.ssoClient.Validate(r.Context(), token)
		// if err != nil {
		// 	http.Error(w, "unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		//ctx := context.WithValue(r.Context(), "userID", resp.UserId)
		ctx := context.WithValue(r.Context(), "userID", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	const prefix = "Bearer "
	if len(authHeader) < len(prefix) {
		return ""
	}

	return authHeader[len(prefix):]
}
