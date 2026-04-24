package app

import (
	"context"
	"fmt"
	deliveryhttp "my-go-server/internal/delivery/http"
	"my-go-server/internal/delivery/http/handler"
	repositorytest "my-go-server/internal/repository/test"
	usecasetest "my-go-server/internal/usecase/test"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	repo := repositorytest.NewRepository()
	service := usecasetest.NewService(repo)
	handler := handler.NewHandler(service)

	srv := http.Server{
		Addr: ":8080",
		Handler: deliveryhttp.SetupRoutes(handler),
	}


	go func() {
		fmt.Println("Listening on :8080")
		err := srv.ListenAndServe()
		if err!=nil {
			fmt.Printf("Stopped listening: %v\n", err)
		}
	} ()
	
	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-shutdown.Done()

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown with error: %v", err)
	}
	fmt.Printf("Shutdown complete.")

}