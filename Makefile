MIGRATIONS_DIR=migrations
DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

build-docker:
	docker compose up -d

build-app:
	go build -o ./bin/api ./cmd/api/main.go

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

migrate-reset:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" reset 