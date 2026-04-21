package app

import (
	"net/http"
	deliveryhttp "my-go-server/internal/delivery/http"
	"my-go-server/internal/delivery/http/handler"
	repositorytest "my-go-server/internal/repository/test"
	usecasetest "my-go-server/internal/usecase/test"
)

func Run() {
	repo := repositorytest.NewRepository()
	service := usecasetest.NewService(repo)
	handler := handler.NewHandler(service)

	deliveryhttp.SetupRoutes(handler)

	err := http.ListenAndServe(":8080", nil)
	if err!=nil {
		panic(err)
	}
}