.PHONY: \
	migrate_create_users_table \
	migrate_create_user_categories_table \
	migrate_create_user_category_relations_table \
	migrate_create_user_contact_relations_table \
	migrate_up_user_all \
	migrate_down_user_all

# Criar migrações
migrate_create_users_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_users_table

migrate_create_user_categories_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_user_categories_table

migrate_create_user_category_relations_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_user_category_relations_table

migrate_create_user_contact_relations_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_user_contact_relations_table

# Executar todas as migrações relacionadas a usuário
migrate_up_user_all:
	@echo "Aplicando todas as migrações: users, user_categories, user_category_relations, user_contact_relations..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

# Reverter todas as migrações relacionadas a usuário
migrate_down_user_all:
	@echo "Revertendo todas as migrações: users, user_categories, user_category_relations, user_contact_relations..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
