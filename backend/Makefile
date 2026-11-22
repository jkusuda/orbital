-include .env
export

run:
	go run ./cmd/api

# Go vs C++ engine. Start the engine first; ENGINE_ADDR comes from .env.
bench:
	go test ./internal/engine/ -bench . -run '^$$' -benchmem

migrate-up:
	go tool goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	go tool goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	go tool goose -dir migrations postgres "$(DATABASE_URL)" status
