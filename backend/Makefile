-include .env
export

run:
	go run ./cmd/api

migrate-up:
	go tool goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	go tool goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	go tool goose -dir migrations postgres "$(DATABASE_URL)" status
