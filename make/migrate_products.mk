.PHONY: migrate_create_products_table migrate_create_product_categories_table migrate_create_product_category_relations_table

migrate_create_products_table:
	@migrate create -ext sql -dir db/migrations -seq create_products_table

migrate_create_product_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_product_categories_table

migrate_create_product_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_product_category_relations_table
