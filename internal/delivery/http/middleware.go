package deliveryhttp

import (
	"context"
	"my-go-server/internal/jwt"
	"net/http"
	"strings"

	"google.golang.org/genproto/googleapis/cloud/retail/v2"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		next(w, r.WithContext(ctx))
	}
}