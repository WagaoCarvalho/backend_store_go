.PHONY: migrate_create_addresses_table migrate_up_address migrate_down_address

migrate_create_addresses_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_addresses_table

migrate_up_address:
	@echo "Aplicando migrações: address..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

migrate_down_address:
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
