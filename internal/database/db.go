package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(user, password, host, port, name string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, name)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return db, nil
}

func Ping(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	return nil
}
