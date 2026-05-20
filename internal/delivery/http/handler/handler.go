package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"my-go-server/internal/domain"
	"my-go-server/internal/jwt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) domain.Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) Handle() http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request)  {
		if r.Method != http.MethodGet {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return 
		}

		if err := someWork(r.Context()); err != nil {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}

		msg, err := h.uc.GetMessage(r.Context())
		if err != nil {
			http.Error(w, "Failed to get message", http.StatusInternalServerError)
			return 
		}

		fmt.Fprint(w, msg)
	}
}

func (h *Handler) HandleDBTest() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request Body", http.StatusBadRequest)
			return
		}

		value := string(body)

		if value == "" {
			http.Error(w, "Empty body", http.StatusBadRequest)
			return
		}

		if _, err := h.uc.SaveMessage(r.Context(), value); err!=nil {
			http.Error(w, "Failed to save value", http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("saved"))
	}
}

func someWork(ctx context.Context) error {
	select {
	case <-time.After(5 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *Handler) RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		var req struct {
			Login 		string `json:"login"`
			Password 	string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, err := h.uc.RegisterUser(r.Context(), req.Login, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return 
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return 
		}
	}
}

func (h *Handler) LoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		var req struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, err := h.uc.GetUserByLogin(r.Context(), req.Login)
		if err != nil {
			http.Error(w, "failed to get user", http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "invalid login or password", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "invalid login or password", http.StatusUnauthorized)
			return
		}

		token, err := jwt.GenerateToken(user.ID, user.Login)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		session, err := h.uc.CreateSession(r.Context(), user.ID)
		if err != nil {
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		response := struct {
			Token     string `json:"token"`
			SessionID string `json:"session_id,omitempty"`
		}{
			Token:     token,
			SessionID: session.SessionID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}