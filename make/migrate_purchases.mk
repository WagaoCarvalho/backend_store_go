.PHONY: migrate_create_purchases_table migrate_create_purchase_categories_table migrate_create_purchase_category_relations_table migrate_create_purchase_product_relations_table migrate_create_purchase_services_table

migrate_create_purchases_table:
	@migrate create -ext sql -dir db/migrations -seq create_purchases_table

migrate_create_purchase_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_purchase_categories_table

migrate_create_purchase_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_purchase_category_relations_table

migrate_create_purchase_product_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_purchase_product_relations_table

migrate_create_purchase_services_table:
	@migrate create -ext sql -dir db/migrations -seq create_purchase_services_table
