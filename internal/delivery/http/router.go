package deliveryhttp

import (
	"net/http"
	"my-go-server/internal/domain"
)

func SetupRoutes(handler domain.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", handler.Handle())
	mux.HandleFunc("/dbtest", handler.HandleDBTest())
	mux.HandleFunc("/auth/register", handler.RegisterUser())
	mux.HandleFunc("/auth/login", handler.LoginUser())
	mux.HandleFunc("/orders/create", AuthMiddleware(handler.AddNewOrder()))
	mux.HandleFunc("/orders/list", AuthMiddleware(handler.GetOrdersList()))
	return mux
}