include .env

up:
	docker-compose --file docker-compose.yml up   --remove-orphans 

down:
	docker-compose down --remove-orphans

api:
	gow run cmd/sil-api/main.go -e .env 

migration:
	goose -dir internal/db/migrations create $(name) sql

migrate:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} up

rollback:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} down
