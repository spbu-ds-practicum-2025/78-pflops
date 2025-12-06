package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Connect() *pgxpool.Pool {
	_ = godotenv.Load()

	dsn := os.Getenv("AD_DB_DSN")
	if dsn == "" {
		dsn = "postgres://ad_user:ad_pass@localhost:5432/ad_db?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to ad database: %v", err)
	}
	fmt.Println("Connected to Ad PostgreSQL")
	return pool
}
