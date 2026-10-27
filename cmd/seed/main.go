package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	ctx := context.Background()

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("Database url (env: DB_URL) is required")
	}

	db, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	if err := seedUsers(ctx, db); err != nil {
		log.Fatal(err)
	}
}

func seedUsers(ctx context.Context, db *pgxpool.Pool) error {
	file, err := os.Open("seeds/users.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(record) != 6 {
			return fmt.Errorf("expected 6 columns, got %d", len(record))
		}

		query := `
			INSERT INTO users (first_name, last_name, full_name, role, email, department)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (email) DO NOTHING
		`
		_, err = db.Exec(ctx, query, record[0], record[1], record[2], record[3], record[4], record[5])
		if err != nil {
			return fmt.Errorf("insert student %q: %w", record[2], err)
		}
	}

	log.Println("Users: seed data inserted")

	return nil
}
