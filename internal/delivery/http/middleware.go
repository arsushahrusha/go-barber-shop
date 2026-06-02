package deliveryhttp

import (
	"context"
	"my-go-server/internal/contextkeys"
	"my-go-server/internal/domain"
	"my-go-server/internal/jwt"
	"my-go-server/internal/response"
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

func (m *MiddleWareManager) JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		claims, err := jwt.ValidateToken(parts[1])
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "invalid token")
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
			response.Error(w, http.StatusUnauthorized, "user not authenticated")
			return
		}

		session, err := m.uc.GetSessionByUserID(r.Context(), userID)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "failed to get session")
			return
		}

		if session == nil {
			response.Error(w, http.StatusUnauthorized, "session not found")
			return
		}

		if time.Now().After(session.ExpiresAt) {
			response.Error(w, http.StatusUnauthorized, "session expired, please log in again")
			return
		}

		if err := m.uc.UpdateSessionExpiry(r.Context(), session.SessionID); err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to update session")
			return
		}

		next(w, r)
	}
}