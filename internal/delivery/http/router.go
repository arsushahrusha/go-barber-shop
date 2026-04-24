package deliveryhttp

import (
	"net/http"
	"my-go-server/internal/delivery/http/handler"
)

func SetupRoutes(handler *handler.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", handler.Handle)
	return mux
}