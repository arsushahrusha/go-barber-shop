package deliveryhttp

import (
	"net/http"
	"my-go-server/internal/domain"
)

func SetupRoutes(handler domain.Handler, uc domain.UseCase) http.Handler {
	mux := http.NewServeMux()

	middlewareManager := NewMiddlewareManager(uc)

	mux.HandleFunc("/health", handler.Health())

	mux.HandleFunc("/test", handler.Handle())
	mux.HandleFunc("/dbtest", handler.HandleDBTest())
	mux.HandleFunc("/auth/register", handler.RegisterUser())
	mux.HandleFunc("/auth/login", handler.LoginUser())

	mux.HandleFunc(
		"/orders/create", 
		middlewareManager.JWTMiddleware(
			middlewareManager.SessionMiddleware(
				handler.AddNewOrder(),
			),
		),
	)
	mux.HandleFunc(
		"/orders/list", 
		middlewareManager.JWTMiddleware(
			middlewareManager.SessionMiddleware(
				handler.GetOrdersList(),
				),
			),
		)


	return mux
}