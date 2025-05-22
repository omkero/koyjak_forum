package config

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Pool *pgxpool.Pool
	once sync.Once
)

// Initialize the database pool once
func InitDB() {
	once.Do(func() {
		var err error
		databaseURL := os.Getenv("DATABASE_URL") // Fetch DB URL from env
		if databaseURL == "" {
			log.Fatal("DATABASE_URL is not set in environment variables")
		}

		// warning don use raw postgresql path like this use .env
		Pool, err = pgxpool.New(context.Background(), "postgresql://postgres:root@localhost:5432/koyjak?sslmode=disable")
		if err != nil {
			log.Fatalf("Failed to initialize database connection pool: %v", err)
		}

		log.Println("Database connection pool initialized")
	})
}
