.PHONY: migrate_create_clients_table migrate_create_client_categories_table migrate_create_client_category_relations_table

migrate_create_clients_table:
	@migrate create -ext sql -dir db/migrations -seq create_clients_table

migrate_create_client_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_client_categories_table

migrate_create_client_category_relations_table:
	@migrate
