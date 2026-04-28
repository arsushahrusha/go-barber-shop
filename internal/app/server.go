package app

import (
	"context"
	"fmt"
	"log"
	"my-go-server/internal/config"
	deliveryhttp "my-go-server/internal/delivery/http"
	"my-go-server/internal/delivery/http/handler"
	"my-go-server/internal/repository"
	database "my-go-server/internal/repository/db"
	usecasetest "my-go-server/internal/usecase/test"
	usecasedb "my-go-server/internal/usecase/db"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := database.NewPostgresDB(config.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	dbrepo := database.NewDBRepository(db)
	
	if err := dbrepo.InitTable(); err != nil {
		log.Fatalf("error initializing test table: %s", err.Error())
	}
	dbservice := usecasedb.NewDBService(dbrepo)
	defer db.Close()

	repo := repository.NewRepository()
	service := usecasetest.NewService(repo)
	handler := handler.NewHandler(service, dbservice)
	 

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