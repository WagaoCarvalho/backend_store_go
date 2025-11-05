.PHONY: migrate_create_client_credits_table  migrate_up_client_credits migrate_down_client_credits

migrate_create_client_credits_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_client_credits_table

migrate_up_client_credits:
	@echo "Aplicando migrações: client_credits..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

migrate_down_client_credits:
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
