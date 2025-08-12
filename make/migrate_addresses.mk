.PHONY: migrate_create_addresses_table migrate_up_address migrate_down_address

migrate_create_addresses_table:
	@migrate create -ext sql -dir db/migrations -seq create_addresses_table

migrate_up_address:
	@echo "Aplicando migrações: address..."
	@migrate -database ${DB_CONN_URL} -path db/migrations up

migrate_down_address:
	@migrate -database ${DB_CONN_URL} -path db/migrations down
