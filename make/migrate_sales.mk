.PHONY: migrate_create_sales_table migrate_create_sale_categories_table migrate_create_sale_category_relations_table migrate_create_sale_product_relations_table migrate_create_sale_services_table

migrate_create_sales_table:
	@migrate create -ext sql -dir db/migrations -seq create_sales_table

migrate_create_sale_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_sale_categories_table

migrate_create_sale_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_sale_category_relations_table

migrate_create_sale_product_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_sale_product_relations_table

migrate_create_sale_services_table:
	@migrate create -ext sql -dir db/migrations -seq create_sale_services_table
