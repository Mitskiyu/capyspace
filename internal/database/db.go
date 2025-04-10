package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	name     = os.Getenv("DB_NAME")
)

func Connect() *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, name)

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
