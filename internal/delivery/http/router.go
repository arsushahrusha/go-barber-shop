package deliveryhttp

import (
	"net/http"
	"my-go-server/internal/delivery/http/handler"
)

func SetupRoutes(handler *handler.Handler) {
	http.HandleFunc("/test", handler.Handle)
}