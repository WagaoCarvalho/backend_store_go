# Carrega variáveis do .env
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: server db migrate_create_users_table migrate_create_wallets_table migrate_create_transactions_table migrate_up migrate_douwn

server:
	@go run cmd/http/*.go

db:
	@docker compose --env-file .env up -d

migrate_create_users_table:
	@migrate create -ext sql -dir db/migrations -seq create_users_table

migrate_create_wallets_table:
	@migrate create -ext sql -dir db/migrations -seq create_wallets_table

migrate_create_transactions_table:
	@migrate create -ext sql -dir db/migrations -seq create_transactions_table

migrate_up:
	@migrate -database ${DB_CONN_URL} -path db/migrations up

migrate_douwn:
	@migrate -database ${DB_CONN_URL} -path db/migrations down

print_env:
	@echo "DB_USER=${DB_USER}"
	@echo "DB_PASSWORD=${DB_PASSWORD}"
	@echo "DB_NAME=${DB_NAME}"
	@echo "DB_HOST=${DB_HOST}"
	@echo "DB_PORT=${DB_PORT}"
	@echo "JWT_SECRET=${JWT_SECRET}"
