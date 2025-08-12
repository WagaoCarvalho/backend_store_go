.PHONY: migrate_create_suppliers_table migrate_create_supplier_categories_table migrate_create_supplier_category_relations_table migrate_up_supplier migrate_down_supplier

migrate_create_suppliers_table:
	@migrate create -ext sql -dir db/migrations -seq create_suppliers_table

migrate_create_supplier_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_supplier_categories_table

migrate_create_supplier_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_supplier_category_relations_table

migrate_down_supplier:
	@migrate -database ${DB_CONN_URL} -path db/migrations down

migrate_up_supplier:
	@echo "Aplicando migrações: supplier..."
	@migrate -database ${DB_CONN_URL} -path db/migrations up
