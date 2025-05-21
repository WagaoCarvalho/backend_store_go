.PHONY: migrate_create_users_table migrate_create_users_categories_table migrate_create_user_category_relations_table

migrate_create_users_table:
	@migrate create -ext sql -dir db/migrations -seq create_users_table

migrate_create_users_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_user_categories_table

migrate_create_user_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_user_category_relations_table
