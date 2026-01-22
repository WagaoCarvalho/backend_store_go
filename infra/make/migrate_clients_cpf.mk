.PHONY: \
	migrate_create_clients_cpf_table \
	migrate_create_client_cpf_contact_relations_table \
	migrate_up_clients_cpf_all \
	migrate_down_clients_cpf_all

# Criação de migrations
migrate_create_clients_cpf_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_clients_cpf_table

migrate_create_client_cpf_contact_relations_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_client_cpf_contact_relations_table

# Aplicar todas as migrations relacionadas a clients_cpf
migrate_up_client_cpf_all:
	@echo "Aplicando todas as migrações: clients_cpf e client_cpf_contact_relations..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

# Reverter todas as migrations relacionadas a clients_cpf
migrate_down_client_cpf_all:
	@echo "Revertendo todas as migrações: clients_cpf e client_cpf_contact_relations..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
