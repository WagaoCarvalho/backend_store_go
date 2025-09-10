.PHONY: migrate_create_sale_items_table migrate_up_sale_items migrate_down_sale_items

migrate_create_sale_items_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_sale_items_table

migrate_up_sale_items:
	@echo "Aplicando migrações: sale_items..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

migrate_down_sale_items:
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
