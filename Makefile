.PHONY: run build test clean migrate-create docker-up docker-down docker-build seed
-include .env 
export
CONNECTION_STRING=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

run:
	go run ./cmd/api/

build:
	go build -o bin/api ./cmd/api/


test:
	go test ./...

clean:
	rm -rf bin/

migrate-create:
	migrate create -ext sql -dir cmd/migrate/migrations -seq $(name)

migrate-up:
	migrate -path cmd/migrate/migrations -database "$(CONNECTION_STRING)" up

migrate-down:
	migrate -path cmd/migrate/migrations -database "$(CONNECTION_STRING)" down 1

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

seed:
	go run ./cmd/migrate/seed/