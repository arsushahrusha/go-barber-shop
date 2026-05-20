package config

import (
	"fmt"
	"os"
)

type Config struct {
	DB DBConfig
	Server ServerConfig
}

type DBConfig struct {
	Host string
	Port string
	Username string
	Password string
	DBName string
	SSLMode string
}

type ServerConfig struct {
	Port string
}

func New() (*Config, error) {
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