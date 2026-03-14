.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

run:
	go run ./cmd/main.go

start:
	docker compose up -d

stop:
	docker compose down -v

migrate:
	docker exec -i shortener-postgres psql -U shortener -d shortener < migrations/001_init.sql