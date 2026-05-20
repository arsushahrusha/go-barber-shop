package db

import (
	"context"
	"fmt"
	"my-go-server/internal/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DSN())
	
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}