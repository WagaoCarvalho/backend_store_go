.PHONY: \
	migrate_create_contacts_table \
	migrate_up_contacts \
	migrate_down_contacts

# Criação de migration
migrate_create_contacts_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_contacts_table

# Aplicar migration de contacts
migrate_up_contact:
	@echo "Aplicando migração: contacts..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

# Reverter migration de contacts
migrate_down_contact:
	@echo "Desfazendo migração: contacts..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
