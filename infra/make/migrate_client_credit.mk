.PHONY: migrate_create_client_credit_table  migrate_up_client_credit migrate_down_client_credit

migrate_create_client_credit_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_client_credit_table

migrate_up_client_credit:
	@echo "Aplicando migrações: client_credit..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

migrate_down_client_credit:
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
