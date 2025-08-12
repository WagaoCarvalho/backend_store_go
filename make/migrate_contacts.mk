.PHONY: migrate_create_contacts_table migrate_up_contact migrate_down_contact

migrate_create_contacts_table:
	@migrate create -ext sql -dir db/migrations -seq create_contacts_table

migrate_up_contact:
	@echo "Aplicando migrações: contact..."
	@migrate -database ${DB_CONN_URL} -path db/migrations up

migrate_down_contact:
	@migrate -database ${DB_CONN_URL} -path db/migrations down

