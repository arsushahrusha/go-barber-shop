package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DBRepository struct {
	db *sqlx.DB
}

func NewDBRepository(db *sqlx.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (r *DBRepository) Save(value string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO db_test (value) values ($1) RETURNING id")
	
	row := r.db.QueryRow(query, value)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DBRepository) InitTable() error {
	query := `CREATE TABLE IF NOT EXISTS db_test (
		id SERIAL PRIMARY KEY,
		value TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err := r.db.Exec(query)
	return err
}