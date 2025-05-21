.PHONY: migrate_create_addresses_table

migrate_create_addresses_table:
	@migrate create -ext sql -dir db/migrations -seq create_addresses_table
