package handler

import (
	"fmt"
	"io"
	"my-go-server/internal/domain"
	"net/http"
	"time"
	"context"
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