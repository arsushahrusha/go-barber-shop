package deliveryhttp

import (
	"net/http"
	"my-go-server/internal/domain"
)

func SetupRoutes(handler domain.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", handler.Handle())
	mux.HandleFunc("/dbtest", handler.HandleDBTest())
	return mux
}