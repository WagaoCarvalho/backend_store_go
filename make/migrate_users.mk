
.PHONY: migrate_create_users_table migrate_create_users_categories_table migrate_create_user_category_relations_table migrate_up_user migrate_down_user

migrate_create_users_table:
	@migrate create -ext sql -dir db/migrations -seq create_users_table

migrate_create_users_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_user_categories_table

migrate_create_user_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_user_category_relations_table

migrate_up_user:
	@echo "Aplicando migrações: user..."
	@migrate -database ${DB_CONN_URL} -path db/migrations up

migrate_down_user:
	@migrate -database ${DB_CONN_URL} -path db/migrations down


