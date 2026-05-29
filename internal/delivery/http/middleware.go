package deliveryhttp

import (
	"context"
	"my-go-server/internal/contextkeys"
	"my-go-server/internal/domain"
	"my-go-server/internal/jwt"
	"net/http"
	"strings"
	"time"
)

type MiddleWareManager struct {
	uc domain.UseCase
}

func NewMiddlewareManager(uc domain.UseCase) *MiddleWareManager {
	return &MiddleWareManager{uc: uc}
}

func (m *MiddleWareManager) JWTMiddleware (next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextkeys.UserID, claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

func (m *MiddleWareManager) SessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(contextkeys.UserID).(string)
		if !ok || userID == "" {
			http.Error(w, "user not authenticated", http.StatusUnauthorized)
			return 
		}

		session, err := m.uc.GetSessionByUserID(r.Context(), userID)
		if err != nil {
			http.Error(w, "failed to get session", http.StatusUnauthorized)
			return 
		}

		if time.Now().After(session.ExpiresAt) {
			http.Error(w, "session expired, please log in again", http.StatusUnauthorized)
			return 
		}

		if err := m.uc.UpdateSessionExpiry(r.Context(), session.SessionID); err != nil {
			http.Error(w, "failed to update session", http.StatusInternalServerError)
			return 
		}

		next(w, r)
	}
}