.PHONY: migrate_create_clients_table  migrate_up_clients migrate_down_clients

migrate_create_clients_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_clients_table

migrate_up_client:
	@echo "Aplicando migrações: client..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

migrate_down_client:
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
