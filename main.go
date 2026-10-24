package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theonlysroy/attendai-be-go/internal/http"
	"github.com/theonlysroy/attendai-be-go/internal/storage/postgres"
)

func main() {

	ctx := context.Background()

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "postgres://postgres:attendai@localhost:5432/attendai"
	}
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
