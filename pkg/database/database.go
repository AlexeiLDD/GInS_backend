package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	connStr := "user=postgres password=postgres dbname=HiveBuddy host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
		return nil, err
	}
	log.Println("Connected to the database successfully")
	return db, nil
}
