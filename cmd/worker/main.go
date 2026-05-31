package main

import (
	"context"
	"fmt"
	"log"
	"my-go-server/internal/config"
	"my-go-server/internal/rabbitmq"
	workerrepo "my-go-server/internal/worker_service/repository"
	workerservice "my-go-server/internal/worker_service/service"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	rabbitClient, err := rabbitmq.New(cfg.RabbitMQ.URL)

	if err != nil {
		log.Fatalf("failed to connect rabbitmq: %s", err.Error())
	}
	defer rabbitClient.Close()

	if err := rabbitClient.DeclareQueue(cfg.RabbitMQ.NewOrderQueue); err != nil {
		log.Fatalf("failed to declare new order queue: %s", err.Error())
	}

	if err := rabbitClient.DeclareQueue(cfg.RabbitMQ.StatusQueue); err != nil {
		log.Fatalf("failed to declare status queue: %s", err.Error())
	}

	repo := workerrepo.NewRabbitRepository(
		rabbitClient,
		cfg.RabbitMQ.NewOrderQueue,
		cfg.RabbitMQ.StatusQueue,
	)

	workerPool := workerservice.NewWorkerPool(repo, cfg.Worker.WorkersCount)

	fmt.Println("worker service started...")
	workerPool.Run(ctx)
	fmt.Println("worker service stopped")
}