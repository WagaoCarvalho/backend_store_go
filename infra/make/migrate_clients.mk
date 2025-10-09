.PHONY: \
	migrate_create_clients_table \
	migrate_create_client_contact_relations_table \
	migrate_up_clients_all \
	migrate_down_clients_all

# Criação de migrations
migrate_create_clients_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_clients_table

migrate_create_client_contact_relations_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_client_contact_relations_table

# Aplicar todas as migrations relacionadas a clients
migrate_up_client_all:
	@echo "Aplicando todas as migrações: clients e client_contact_relations..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

# Reverter todas as migrations relacionadas a clients
migrate_down_client_all:
	@echo "Revertendo todas as migrações: clients e client_contact_relations..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
