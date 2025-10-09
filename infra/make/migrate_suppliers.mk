.PHONY: \
	migrate_create_suppliers_table \
	migrate_create_supplier_categories_table \
	migrate_create_supplier_category_relations_table \
	migrate_create_supplier_contact_relations_table \
	migrate_up_supplier_all \
	migrate_down_supplier_all

# Criação de migrations
migrate_create_suppliers_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_suppliers_table

migrate_create_supplier_categories_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_supplier_categories_table

migrate_create_supplier_category_relations_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_supplier_category_relations_table

migrate_create_supplier_contact_relations_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_contact_suppliers_table

# Aplicar todas as migrations relacionadas a supplier
migrate_up_supplier_all:
	@echo "Aplicando todas as migrações: supplier e contact_supplier..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

# Reverter todas as migrations relacionadas a supplier
migrate_down_supplier_all:
	@echo "Desfazendo todas as migrações: supplier e contact_supplier..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
