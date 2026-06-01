package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"my-go-server/internal/contextkeys"
	"my-go-server/internal/domain"
	"my-go-server/internal/jwt"
	"my-go-server/internal/models"
	"my-go-server/internal/response"
	"net/http"
	"strconv"
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
			response.Error(w, http.StatusMethodNotAllowed, "method is not allowed")
			return
		}

		defer r.Body.Close()

		var req struct {
			Login 		string `json:"login"`
			Password 	string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Login == "" {
			response.Error(w, http.StatusBadRequest, "login is required")
			return
		}

		if req.Password == "" {
			response.Error(w, http.StatusBadRequest, "password is required")
			return
		}

		user, err := h.uc.RegisterUser(r.Context(), req.Login, req.Password)
		if err != nil {
			response.Error(w, http.StatusConflict, err.Error())
			return 
		}

		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusCreated)

		// if err := json.NewEncoder(w).Encode(user); err != nil {
		// 	http.Error(w, "failed to encode response", http.StatusInternalServerError)
		// 	return 
		// }

		response.JSON(w, http.StatusCreated, user)
	}
}

func (h *Handler) LoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "method is not allowed")
			return
		}

		defer r.Body.Close()

		var req struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Login == "" {
			response.Error(w, http.StatusBadRequest, "login is required")
			return
		}

		if req.Password == "" {
			response.Error(w, http.StatusBadRequest, "password is required")
			return
		}

		user, err := h.uc.GetUserByLogin(r.Context(), req.Login)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to get user")
			return
		}

		if user == nil {
			response.Error(w, http.StatusInternalServerError, "failed to get user")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			response.Error(w, http.StatusUnauthorized, "invalid login or password")
			return
		}

		token, err := jwt.GenerateToken(user.ID, user.Login)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to generate token")
			return
		}

		session, err := h.uc.CreateSession(r.Context(), user.ID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to create session")
			return
		}

		response.JSON(w, http.StatusOK, struct {
			Token     string `json:"token"`
			SessionID string `json:"session_id,omitempty"`
		}{
			Token:     token,
			SessionID: session.SessionID,
		})
	}
}

func (h *Handler) AddNewOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "method is not allowed")
			return
		}

		userID, ok := r.Context().Value(contextkeys.UserID).(string)
		if !ok || userID == "" {
			response.Error(w, http.StatusUnauthorized, "user not authenticated")
			return
		}

		defer r.Body.Close()

		var req models.CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Amount <= 0 {
			response.Error(w, http.StatusBadRequest, "amount is required and must be positive")
			return
		}

		orderIDs := make([]string, 0, req.Amount)

		for i := 0; i < req.Amount; i++ {
			order, err := h.uc.CreateOrder(r.Context(), userID)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, "failed to create order")
				return
			}

			orderIDs = append(orderIDs, order.ID)
		}

		response.JSON(w, http.StatusOK, models.CreateOrderResponse{
			OrderIDs: orderIDs,
		})
	}
}

func (h *Handler) GetOrdersList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "method is not allowed")
			return
		}

		userID, ok := r.Context().Value(contextkeys.UserID).(string)
		if !ok || userID == "" {
			response.Error(w, http.StatusUnauthorized, "user not authenticated")
			return
		}

		activeOnly := false
		activeParam := r.URL.Query().Get("active")

		if activeParam != "" {
			value, err := strconv.ParseBool(activeParam)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "invalid active query param")
				return
			}

			activeOnly = value
		}

		orders, err := h.uc.GetOrdersByUserID(r.Context(), userID, activeOnly)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to get orders")
			return
		}

		response.JSON(w, http.StatusOK, models.OrdersListResponse{
			Count: len(orders),
			Items: orders,
		})
	}
}

func (h *Handler) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "method is not allowed")
			return
		}

		response.JSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	}
}