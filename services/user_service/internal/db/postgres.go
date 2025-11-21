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

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        host := os.Getenv("POSTGRES_HOST")
        if host == "" {
            host = "postgres"
        }
        port := os.Getenv("POSTGRES_CONTAINER_PORT")
        if port == "" {
            port = "5432"
        }
        user := os.Getenv("POSTGRES_USER")
        if user == "" {
            user = "admin"
        }
        pass := os.Getenv("POSTGRES_PASSWORD")
        if pass == "" {
            pass = "password"
        }
        db := os.Getenv("POSTGRES_DB")
        if db == "" {
            db = "user_service"
        }
        dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)
    }
    pool, err := pgxpool.New(context.Background(), dsn)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    fmt.Println("Connected to PostgreSQL")
    return pool
}