DB_URL ?= db/tyres.db
MIGRATIONS_DIR := ./db/migrations


.PHONY: run db-up db-down migration-status new-migration db-test

run:
	go run ./cmd/pin-pyre-db-api

db-test:
	go test ./db/tests/...

db-up:
	goose -dir $(MIGRATIONS_DIR) sqlite3 "$(DB_URL)" up

db-down:
	goose -dir $(MIGRATIONS_DIR) sqlite3 "$(DB_URL)" down

migration-status:
	goose -dir $(MIGRATIONS_DIR) sqlite3 "$(DB_URL)" status

new-migration:
	goose -dir $(MIGRATIONS_DIR) sqlite3 "$(DB_URL)" create $(name) sql
