package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"my-go-server/internal/config"
	deliveryhttp "my-go-server/internal/delivery/http"
	"my-go-server/internal/delivery/http/handler"
	"my-go-server/internal/rabbitmq"
	"my-go-server/internal/repository"
	database "my-go-server/internal/repository/db"
	"my-go-server/internal/usecase"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	cfg, err := config.New()

	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	postgresDB, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer postgresDB.Close()

	dbrepo := database.NewDBRepository(postgresDB)

	initCtx, initCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer initCancel()

	if err := dbrepo.InitTables(initCtx); err != nil {
		log.Fatalf("failed to init database tables: %s", err.Error())
	}

	rabbitClient, err := rabbitmq.New(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatalf("failed to connect rabbitmq: %s", err.Error())
	}
	defer rabbitClient.Close()

	if err := rabbitClient.DeclareQueue(cfg.RabbitMQ.NewOrderQueue); err != nil {
		log.Fatalf("failed to declare new orders queue: %s", err.Error())
	}

	if err := rabbitClient.DeclareQueue(cfg.RabbitMQ.StatusQueue); err != nil {
		log.Fatalf("failed to declare status queue: %s", err.Error())
	}


	repo := repository.NewRepository(
		dbrepo,
		rabbitClient,
		cfg.RabbitMQ.NewOrderQueue,
		cfg.RabbitMQ.StatusQueue,
	)
	uc := usecase.NewService(repo)
	h := handler.NewHandler(uc)
	
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	repo.StartStatusChangeConsumer(appCtx, func(orderID, status string) error {
		return uc.ChangeOrderStatus(appCtx, orderID, status)
	})

	srv := http.Server{
		Addr: ":"+cfg.Server.Port,
		Handler: deliveryhttp.SetupRoutes(h, uc),
	}


	go func() {
		fmt.Printf("Listening on :%s\n", cfg.Server.Port)
		
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %s", err.Error())
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
		return
	}
	fmt.Printf("Shutdown complete.")

}