package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/theonlysroy/attendai-be-go/internal/http"
	"github.com/theonlysroy/attendai-be-go/internal/storage/postgres"
	"github.com/theonlysroy/attendai-be-go/internal/utils"
)

func main() {

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	dsn := "postgres://postgres:attendai@localhost:5432/attendaidb"
	dbUrl := utils.Getkey("DB_URL", dsn)

	db, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	studentStore := postgres.NewStudentStore(db)

	server := http.NewAPIServer(":4040", studentStore)
	server.Run()
}
