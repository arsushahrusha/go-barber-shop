package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DB DBConfig
	Server ServerConfig
	RabbitMQ RabbitMQConfig
	Worker WorkerConfig
}

type DBConfig struct {
	Host string
	Port string
	Username string
	Password string
	DBName string
	SSLMode string
}

type RabbitMQConfig struct {
	URL string
	NewOrderQueue string
	StatusQueue string
}

type WorkerConfig struct {
	WorkersCount int
}

type ServerConfig struct {
	Port string
}

func New() (*Config, error) {
	workerCount, err := strconv.Atoi(os.Getenv("WORKERS_COUNT"))
	if err != nil || workerCount <= 0 {
		workerCount = 3 //by default
	}
	return &Config{
		DB: DBConfig{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName: os.Getenv("DB_NAME"),
			SSLMode: os.Getenv("DB_SSLMODE"),
		},
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: os.Getenv("RABBITMQ_URL"),
			NewOrderQueue: os.Getenv("RABBITMQ_NEW_ORDERS_QUEUE"),
			StatusQueue: os.Getenv("RABBITMQ_STATUS_QUEUE"),
		},
		Worker: WorkerConfig{
			WorkersCount: workerCount,
		},
	}, nil
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Username,
		c.DBName,
		c.Password,
		c.SSLMode,
	)
}