.PHONY: tidy build run test migrate-up migrate-down migrate-status migrate-create compose-up compose-up

migration_dir := "migrations"
db_driver := "postgres"
db_url := "postgres://postgres:attendai@localhost:5432/attendai?sslmode=disable"

tidy:
	@go mod tidy

build:
	@go build -o bin/attendai .

run: build
	@./bin/attendai

test:
	@go test -v ./...

migrate-up:
	@goose -dir $(migration_dir) $(db_driver) "$(db_url)" up

migrate-down:
	@goose -dir $(migration_dir) $(db_driver) "$(db_url)" down

migrate-status:
	@goose -dir $(migration_dir) $(db_driver) "$(db_url)" status

migrate-create:
	@goose -dir $(migration_dir) create $(name) sql

compose-up:
	@docker compose up -d

compose-down:
	@docker compose down -v
