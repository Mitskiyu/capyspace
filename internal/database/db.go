package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(connStr string) *sql.DB {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not ping db: %v", err)
	}

	log.Println("Connected to database!")
	return db
}
