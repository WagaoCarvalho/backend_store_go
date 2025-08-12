.PHONY: migrate_create_products_table migrate_create_product_categories_table migrate_create_product_category_relations_table migrate_up_product migrate_down_product

migrate_create_products_table:
	@migrate create -ext sql -dir db/migrations -seq create_products_table

migrate_create_product_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_product_categories_table

migrate_create_product_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_product_category_relations_table

migrate_up_product:
	@echo "Aplicando migrações: product..."
	@migrate -database ${DB_CONN_URL} -path db/migrations up

migrate_down_product:
	@migrate -database ${DB_CONN_URL} -path db/migrations down
