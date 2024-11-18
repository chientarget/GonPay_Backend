package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDB(connStr string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &PostgresDB{DB: db}, nil
}
